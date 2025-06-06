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

package main

import (
	"text/template"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/metric"
	"github.com/vdaas/vald/internal/observability/metrics/backoff"
	"github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/observability/metrics/grpc"
	"github.com/vdaas/vald/internal/observability/metrics/version"
)

const (
	heightShort  = 3
	heightMedium = 6
	heightTall   = 8

	widthOneEighth = 3
	widthOneSixth  = 4
	widthQuarter   = 6
	witdhOneThird  = 8
	widthHalf      = 12
	widthFull      = 24

	opacity = 10

	organization      = "vdaas"
	repository        = "vald"
	defaultMaintainer = organization + ".org " + repository + " team <" + repository + "@" + organization + ".org>"
	header            = `#
# Copyright (C) 2019-{{.Year}} {{.Maintainer}}
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#`
	maintainerKey = "MAINTAINER"
	rootKey       = "ROOTDIR"

	namespaceKey        = "exported_kubernetes_namespace"
	namespaceVariable   = "$Namespace"
	nameKey             = "kubernetes_name"
	appNameKey          = "app_kubernetes_io_name"
	nameVariable        = "$ReplicaSet"
	podKey              = "target_pod"
	instanceKey         = "instance"
	podVariable         = "$PodName"
	imageKey            = "image"
	cpuMetric           = "container_cpu_usage_seconds_total"
	nodeCPUMetric       = "node_cpu_seconds_total"
	memMetric           = "container_memory_working_set_bytes"
	intervalVariable    = "$interval"
	reconcileSec        = "controller_runtime_reconcile_time_seconds_bucket"
	reconcileTotal      = "controller_runtime_reconcile_total"
	reconcileErrorTotal = "controller_runtime_reconcile_errors_total"
	restDurationSec     = "rest_client_request_duration_seconds_bucket"
	helmInfo            = "helm_operator_build_info"
	goInfo              = "go_info"
	resCreatedSec       = "resource_created_at_seconds"
	backoffRetryCount   = backoff.MetricsName
	serverCompletedRPCs = grpc.CompletedRPCsMetricsName
	podInfo             = "kube_pod_info"
	appInfo             = version.Name
	grpcServerMethod    = metric.GRPCMethodKeyName
	grpcServerStatus    = metric.GRPCStatus
	serverLatencyBucket = "server_latency_bucket"
	circuitBreakerState = circuitbreaker.MetricsName
)

type kindStatus struct {
	kind   string
	status string
}

var (
	allKindStatus = []kindStatus{
		{"deployment", "replicas"},
		{"deployment", "replicas_available"},
		{"deployment", "replicas_unavailable"},
		{"deployment", "replicas_updated"},
		{"daemonset", "current_number_scheduled"},
		{"daemonset", "desired_number_scheduled"},
		{"daemonset", "number_available"},
		{"daemonset", "number_misscheduled"},
		{"daemonset", "number_ready"},
		{"daemonset", "number_unavailable"},
		{"statefulset", "replicas_current"},
		{"statefulset", "replicas_ready"},
		{"statefulset", "replicas_updated"},
	}

	license = template.Must(template.New("license").Parse(header + `

# DO_NOT_EDIT this workflow file is generated by https://github.com/vdaas/vald/blob/main/hack/grafana/gen/src/main.go
`))

	quntiles = []float64{0.5, 0.95, 0.99}

	indiceThresholds = []dashboard.Threshold{
		{Value: nil, Color: "red"},
		{Value: cog.ToPtr[float64](100000), Color: "orange"},
		{Value: cog.ToPtr[float64](10000000), Color: "green"},
	}
	podThresholds = []dashboard.Threshold{
		{Value: nil, Color: "green"},
		{Value: cog.ToPtr[float64](1), Color: "orange"},
		{Value: cog.ToPtr[float64](5), Color: "red"},
	}
	queueThresholds = []dashboard.Threshold{
		{Value: nil, Color: "green"},
		{Value: cog.ToPtr[float64](100), Color: "orange"},
		{Value: cog.ToPtr[float64](1000), Color: "red"},
	}
	memThresholds = []dashboard.Threshold{
		{Value: nil, Color: "green"},
		{Value: cog.ToPtr[float64](10000000000), Color: "orange"},
		{Value: cog.ToPtr[float64](1000000000000), Color: "red"},
	}
)
