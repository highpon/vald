// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"os"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
)

// Loader is representation of load test.
type Loader interface {
	Prepare(context.Context) error
	Do(context.Context) <-chan error
}

type (
	loadFunc func(context.Context, *grpc.ClientConn, any, ...grpc.CallOption) (any, error)
)

type loader struct {
	eg               errgroup.Group
	client           grpc.Client
	addr             string
	concurrency      int
	batchSize        int
	dataset          string
	progressDuration time.Duration
	loaderFunc       loadFunc
	sendDataProvider func() *any
	dataSize         int
	operation        config.Operation
}

// NewLoader returns Loader implementation.
func NewLoader(opts ...Option) (Loader, error) {
	l := new(loader)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(l); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	var err error
	switch l.operation {
	case config.Insert:
		l.loaderFunc, err = l.newInsert()
	case config.StreamInsert:
		l.loaderFunc, err = l.newStreamInsert()
	case config.Search:
		l.loaderFunc, err = l.newSearch()
	case config.StreamSearch:
		l.loaderFunc, err = l.newStreamSearch()
	default:
		err = errors.Errorf("undefined operation: %s", l.operation.String())
	}
	if err != nil {
		return nil, err
	}

	return l, nil
}

// Prepare generate request data.
func (l *loader) Prepare(context.Context) (err error) {
	fn := assets.Data(l.dataset)
	if fn == nil {
		return errors.Errorf("dataset load function is nil: %s", l.dataset)
	}
	dataset, err := fn()
	if err != nil {
		return err
	}

	switch l.operation {
	case config.Insert, config.StreamInsert:
		l.sendDataProvider, l.dataSize, err = insertRequestProvider(dataset, l.batchSize)
	case config.Search, config.StreamSearch:
		l.sendDataProvider, l.dataSize, err = searchRequestProvider(dataset)
	}
	if err != nil {
		return err
	}

	return nil
}

// Do operates load test.
func (l *loader) Do(ctx context.Context) <-chan error {
	ech := make(chan error, l.dataSize)
	finalize := func(ctx context.Context, err error) {
		select {
		case <-ctx.Done():
			ech <- errors.Join(err, ctx.Err())
		case ech <- err:
		}
	}

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		finalize(ctx, err)
		return ech
	}

	var pgCnt, errCnt int32 = 0, 0
	var start, end time.Time
	vps := func(count int, t1, t2 time.Time) float64 {
		return float64(count) / t2.Sub(t1).Seconds()
	}
	progress := func() {
		log.Infof("progress %d requests, %f[vps], error: %d", pgCnt, vps(int(pgCnt)*l.batchSize, start, time.Now()), errCnt)
	}

	f := func(i *any, err error) {
		atomic.AddInt32(&pgCnt, 1)
		if err != nil {
			atomic.AddInt32(&errCnt, 1)
		}
	}

	ticker := time.NewTicker(l.progressDuration)
	l.eg.Go(safety.RecoverFunc(func() error {
		defer ticker.Stop()
		for int(pgCnt) < l.dataSize {
			if err := func() error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					progress()
					return nil
				}
			}(); err != nil {
				return err
			}
		}
		return nil
	}))

	l.eg.Go(safety.RecoverFunc(func() error {
		log.Infof("start load test(%s)", l.operation.String())
		defer close(ech)
		defer ticker.Stop()
		start = time.Now()
		err := l.do(ctx, f, finalize)
		end = time.Now()

		if errCnt > 0 {
			log.Warnf("Error ratio: %.2f%%", float64(errCnt)/float64(pgCnt)*100)
		}
		if err != nil {
			finalize(ctx, err)
			return p.Signal(syscall.SIGKILL) // TODO: #403
		}
		log.Infof("result:%d\t%d\t%f", l.concurrency, l.batchSize, vps(int(pgCnt)*l.batchSize, start, end))

		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))
	return ech
}

func (l *loader) do(
	ctx context.Context, f func(*any, error), notify func(context.Context, error),
) (err error) {
	eg, egctx := errgroup.New(ctx)

	switch l.operation {
	case config.StreamInsert, config.StreamSearch:
		eg.Go(safety.RecoverFunc(func() (err error) {
			defer func() {
				if err != nil {
					notify(egctx, err)
					err = nil
				}
			}()
			_, err = l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
				st, err := l.loaderFunc(ctx, conn, nil, copts...)
				if err != nil {
					return nil, err
				}

				if l.operation == config.StreamInsert {
					return nil, grpc.BidirectionalStreamClient(st.(grpc.ClientStream), l.sendDataProvider, func(i *payload.Empty, err error) bool {
						f(nil, err)
						return true
					})
				} else {
					return nil, grpc.BidirectionalStreamClient(st.(grpc.ClientStream), l.sendDataProvider, func(i *payload.Search_Response, err error) bool {
						f(nil, err)
						return true
					})
				}
			})
			return err
		}))
		err = eg.Wait()
	case config.Insert, config.Search:
		eg.SetLimit(l.concurrency)

		for {
			r := l.sendDataProvider()
			if r == nil {
				break
			}

			eg.Go(safety.RecoverFunc(func() (err error) {
				defer func() {
					notify(egctx, err)
					err = nil
				}()
				_, err = l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
					res, err := l.loaderFunc(egctx, conn, *r)
					f(&res, err)
					return res, err
				})

				return err
			}))
		}
		err = eg.Wait()
	default:
		err = errors.Errorf("undefined type: %s", l.operation.String())
	}
	return
}
