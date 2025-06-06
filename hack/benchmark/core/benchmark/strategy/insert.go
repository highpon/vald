//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package strategy provides benchmark strategy
package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
)

func NewInsert(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Insert"),
		WithProp32(
			func(ctx context.Context, b *testing.B, c algorithm.Bit32, dataset assets.Dataset, ids []uint, cnt *uint64) (any, error) {
				v, err := dataset.Train(int(atomic.LoadUint64(cnt)) % dataset.TrainSize())
				if err != nil {
					return nil, err
				}
				b.StartTimer()
				defer b.StopTimer()
				return c.Insert(v.([]float32))
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c algorithm.Bit64, dataset assets.Dataset, ids []uint, cnt *uint64) (any, error) {
				v, err := dataset.Train(int(atomic.LoadUint64(cnt)) % dataset.TrainSize())
				if err != nil {
					return nil, err
				}
				b.StartTimer()
				defer b.StopTimer()
				return c.Insert(float32To64(v.([]float32)))
			},
		),
	}, opts...)...)
}
