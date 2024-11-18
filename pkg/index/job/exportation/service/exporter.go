// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	apiName        = "vald/index/job/export"
	grpcMethodName = "vald.v1.StreamListObject/" + vald.StreamListObjectRPCName
)

// Exporter represents an interface for exporting.
type Exporter interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
}

type export struct {
	eg           errgroup.Group
	gateway      vc.Client
	storedVector pogreb.DB

	streamListConcurrency        int
	backgroundSyncInterval       time.Duration
	backgroundCompactionInterval time.Duration
	indexPath                    string
}

// New returns Exporter object if no error occurs.
func New(opts ...Option) (Exporter, error) {
	e := new(export)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(e); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	if err := file.MkdirAll(e.indexPath, os.ModePerm); err != nil {
		log.Errorf("failed to create dir %s", e.indexPath)
		return nil, err
	}
	// Todo: Determine file name
	path := file.Join(e.indexPath, fmt.Sprintf("%s.db", strconv.FormatInt(time.Now().Unix(), 10)))
	db, err := pogreb.New(pogreb.WithPath(path),
		pogreb.WithBackgroundCompactionInterval(e.backgroundCompactionInterval),
		pogreb.WithBackgroundSyncInterval(e.backgroundSyncInterval))
	if err != nil {
		log.Errorf("failed to open checked List kvs DB %s", path)
		return nil, err
	}
	e.storedVector = db
	return e, nil
}

// StartClient starts the gRPC client.
func (e *export) StartClient(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)
	gch, err := e.gateway.Start(ctx)
	if err != nil {
		return nil, err
	}
	e.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-gch:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (e *export) Start(ctx context.Context) error {
	err := e.doExportIndex(ctx,
		func(ctx context.Context, rc vald.ObjectClient, copts ...grpc.CallOption) (vald.Object_StreamListObjectClient, error) {
			return rc.StreamListObject(ctx, &payload.Object_List_Request{}, copts...)
		},
	)
	return err
}

func (e *export) doExportIndex(
	ctx context.Context,
	fn func(ctx context.Context, rc vald.ObjectClient, copts ...grpc.CallOption) (vald.Object_StreamListObjectClient, error),
) (errs error) {
	ctx, span := trace.StartSpan(igrpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doExportIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	emptyReq := new(payload.Object_List_Request)

	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(e.streamListConcurrency)
	ctx, cancel := context.WithCancelCause(egctx)
	gatewayAddrs := e.gateway.GRPCClient().ConnectedAddrs()
	if len(gatewayAddrs) < 0 {
		log.Errorf("Active gateway is not found.: %v ", ctx.Err())
	}

	conn, err := grpc.NewClient(gatewayAddrs[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	vcClient := vc.NewValdClient(conn)
	grpcCallOpts := []grpc.CallOption{
		grpc.WaitForReady(true),
	}
	// stream, err := fn(ctx, vc.NewValdClient(conn), grpcCallOpts...)
	stream, err := vcClient.StreamListObject(ctx, emptyReq, grpcCallOpts...)
	if err != nil || stream == nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				log.Errorf("context done unexpectedly: %v", ctx.Err())
			}
			if !errors.Is(context.Cause(ctx), io.EOF) {
				log.Errorf("context canceled due to %v", ctx.Err())
			}
			err = eg.Wait()
			if err != nil {
				log.Errorf("exporter returned error status errgroup returned error: %v", ctx.Err())
			} else {
				log.Infof("exporter finished")
			}
			return nil
		default:
			res, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					cancel(io.EOF)
				} else {
					cancel(errors.ErrStreamListObjectStreamFinishedUnexpectedly(err))
				}
			} else if res != nil && res.GetVector() != nil && res.GetVector().GetId() != "" {
				eg.Go(safety.RecoverFunc(func() (err error) {
					objVec := res.GetVector()
					log.Infof("received object vector id: %s, timestamp: %d", objVec.GetId(), objVec.GetTimestamp())

					storedBinary, ok, err := e.storedVector.Get(objVec.GetId())
					if err != nil {
						log.Errorf("failed to perform Get from check list but still try to finish processing without cache: %v", err)
						return err
					}

					var storedObjVec payload.Object_Vector
					if ok {
						if err := storedObjVec.UnmarshalVT(storedBinary); err != nil {
							log.Errorf("failed to Unmarshal proto to payload.Object_Vector: %v", err)
							return err
						}
					}

					isUpsertVector := !ok || storedObjVec.GetTimestamp() < objVec.GetTimestamp()
					if isUpsertVector {
						dAtA, err := objVec.MarshalVT()
						if err != nil {
							return err
						}
						e.storedVector.Set(objVec.GetId(), dAtA)
					}
					return nil
				}))
			}
		}
	}
}
