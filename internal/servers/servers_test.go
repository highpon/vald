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
package servers

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		checkFunc func(got, want *listener) error
		want      *listener
	}

	tests := []test{
		{
			name: "initialize with default options",
			want: &listener{
				eg: errgroup.Get(),
			},
			checkFunc: func(got *listener, want *listener) error {
				if !reflect.DeepEqual(got, want) {
					return errors.Errorf("not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
		{
			name: "initialize with custom options",
			opts: []Option{
				WithStartUpStrategy([]string{
					"strg_1",
					"strg_2",
				}),
			},
			want: &listener{
				eg: errgroup.Get(),
				sus: []string{
					"strg_1",
					"strg_2",
				},
				sds: []string{
					"strg_2",
					"strg_1",
				},
			},
			checkFunc: func(got *listener, want *listener) error {
				if !reflect.DeepEqual(got.sus, want.sus) {
					return errors.Errorf("sus is not equals. want: %v, got: %v", want.sus, got.sus)
				}

				if !reflect.DeepEqual(got.sds, want.sds) {
					return errors.Errorf("sds is not equals. want: %v, got: %v", want.sds, got.sds)
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			got := New(test.opts...)
			if err := test.checkFunc(got.(*listener), test.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_listener_ListenAndServe(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		eg      errgroup.Group
		servers map[string]server.Server
		sus     []string
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got, want <-chan error) error
		afterFunc func(*testing.T)
		want      <-chan error
	}

	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			srv1 := &mockServer{
				IsRunningFunc: func() bool {
					return false
				},
				ListenAndServeFunc: func(context.Context, chan<- error) error {
					return nil
				},
			}

			srv2Err := errors.New("srv2 error")
			srv2 := &mockServer{
				IsRunningFunc: func() bool {
					return false
				},
				ListenAndServeFunc: func(context.Context, chan<- error) error {
					return srv2Err
				},
			}

			servers := map[string]server.Server{
				"srv1": srv1,
				"srv2": srv2,
			}

			sus := []string{
				"srv1",
				"srv2",
				"srv3",
			}

			errCh := make(chan error, len(servers)+1)
			errCh <- srv2Err
			errCh <- errors.ErrServerNotFound("srv3")
			errCh <- srv2Err
			close(errCh)

			return test{
				name: "ListenAndServe is success",
				args: args{
					ctx: func() context.Context {
						ctx, cancel := context.WithCancel(ctx)
						go func() {
							defer cancel()
							time.Sleep(time.Second)
						}()
						return ctx
					}(),
				},
				field: field{
					eg:      eg,
					servers: servers,
					sus:     sus,
				},
				want: errCh,
				checkFunc: func(got <-chan error, want <-chan error) error {
					gerrs := make([]error, 0, len(servers))
					for err := range got {
						gerrs = append(gerrs, err)
					}

					werrs := make([]error, 0, len(servers))
					for err := range want {
						werrs = append(werrs, err)
					}

					if len(werrs) != len(gerrs) {
						return errors.Errorf("errors count is not equals: want: %v, got: %v", werrs, gerrs)
					}
					for i := range werrs {
						if gerrs[i].Error() != werrs[i].Error() {
							return errors.Errorf("errors[%d] is not equals: want: %v, got: %v", i, werrs[i], gerrs[i])
						}
					}
					return nil
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			ctx, cancel := context.WithCancel(test.args.ctx)
			defer cancel()
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}

			l := &listener{
				eg:      test.field.eg,
				servers: test.field.servers,
				sus:     test.field.sus,
			}

			errCh := l.ListenAndServe(ctx)
			if err := test.checkFunc(errCh, test.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_listener_Shutdown(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		eg      errgroup.Group
		servers map[string]server.Server
		sds     []string
		sddur   time.Duration
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got, want error) error
		afterFunc func(*testing.T)
		want      error
	}

	defaultCheckFunc := func(got, want error) error {
		if !errors.Is(want, got) {
			return errors.Errorf("not equals. want: %v, got: %v", want, got)
		}
		return nil
	}

	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			srv1 := &mockServer{
				IsRunningFunc: func() bool {
					return true
				},
				ShutdownFunc: func(context.Context) error {
					return nil
				},
			}

			srv2 := &mockServer{
				IsRunningFunc: func() bool {
					return true
				},
				ShutdownFunc: func(context.Context) error {
					return nil
				},
			}

			servers := map[string]server.Server{
				"srv1": srv1,
				"srv2": srv2,
			}

			sds := []string{
				"srv1",
				"srv2",
			}

			return test{
				name: "Shutdown is success",
				args: args{
					ctx: ctx,
				},
				field: field{
					eg:      eg,
					servers: servers,
					sds:     sds,
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			return test{
				name: "server not found error",
				args: args{
					ctx: ctx,
				},
				field: field{
					eg:      eg,
					servers: map[string]server.Server{},
					sds: []string{
						"srv1",
					},
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
				want: errors.ErrServerNotFound("srv1"),
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			want := errors.Wrap(errors.Errorf("unexpected error"), "faild to shutdown")

			srv1 := &mockServer{
				IsRunningFunc: func() bool {
					return true
				},
				ShutdownFunc: func(context.Context) error {
					return want
				},
			}

			servers := map[string]server.Server{
				"srv1": srv1,
			}

			sds := []string{
				"srv1",
			}

			return test{
				name: "unexpected error",
				args: args{
					ctx: ctx,
				},
				field: field{
					eg:      eg,
					servers: servers,
					sds:     sds,
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
				want: want,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			ctx, cancel := context.WithCancel(test.args.ctx)
			defer cancel()
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			l := &listener{
				eg:      test.field.eg,
				servers: test.field.servers,
				sds:     test.field.sds,
				sddur:   test.field.sddur,
			}

			err := l.Shutdown(ctx)
			if err := test.checkFunc(err, test.want); err != nil {
				t.Error(err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
