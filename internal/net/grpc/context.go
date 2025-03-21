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
package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/strings"
)

type contextKey string

// GRPCMethodContextKey represents a context key for gRPC method.
// This is exported only for testing.
const GRPCMethodContextKey contextKey = "grpc_method"

// WrapGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WrapGRPCMethod(ctx context.Context, method string) context.Context {
	m := FromGRPCMethod(ctx)
	if m == "" {
		return context.WithValue(ctx, GRPCMethodContextKey, method)
	}
	if strings.HasSuffix(m, method) {
		return ctx
	}
	return context.WithValue(ctx, GRPCMethodContextKey, m+"/"+method)
}

// WithGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WithGRPCMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, GRPCMethodContextKey, method)
}

// FromGRPCMethod returns the value associated with this context for key (grpcMethodContextKey).
func FromGRPCMethod(ctx context.Context) string {
	if v := ctx.Value(GRPCMethodContextKey); v != nil {
		if method, ok := v.(string); ok {
			return method
		}
	}
	return ""
}
