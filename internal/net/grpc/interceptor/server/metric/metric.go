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
package metric

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
)

const (
	latencyMetricsName       = "server_latency"
	completedRPCsMetricsName = "server_completed_rpcs"

	GRPCMethodKeyName = "grpc_server_method"
	GRPCStatus        = "grpc_server_status"
)

func MetricInterceptors() (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, error) {
	meter := metrics.GetMeter()

	latencyHistogram, err := meter.Float64Histogram(
		latencyMetricsName,
		metrics.WithDescription("Server latency in milliseconds, by method"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create latency metric")
	}

	completedRPCCnt, err := meter.Int64Counter(
		completedRPCsMetricsName,
		metrics.WithDescription("Count of RPCs by method and status"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create completedRPCs metric")
	}

	record := func(ctx context.Context, method string, err error, latency float64) {
		attrs := attributesFromError(method, err)
		latencyHistogram.Record(ctx, latency, metrics.WithAttributes(attrs...))
		completedRPCCnt.Add(ctx, 1, metrics.WithAttributes(attrs...))
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			now := time.Now()
			resp, err = handler(ctx, req)
			elapsedTime := time.Since(now)
			record(ctx, info.FullMethod, err, float64(elapsedTime)/float64(time.Millisecond))
			return resp, err
		}, func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			now := time.Now()
			err = handler(srv, ss)
			elapsedTime := time.Since(now)
			record(ss.Context(), info.FullMethod, err, float64(elapsedTime)/float64(time.Millisecond))
			return err
		}, nil
}

func attributesFromError(method string, err error) []attribute.KeyValue {
	code := codes.OK // default error is success when error is nil
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil {
			code = st.Code()
		}
	}
	return []attribute.KeyValue{
		attribute.String(GRPCMethodKeyName, method),
		attribute.String(GRPCStatus, code.String()),
	}
}
