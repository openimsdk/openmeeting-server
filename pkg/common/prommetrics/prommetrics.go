// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prommetrics

import (
	gp "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"openmeeting-server/pkg/common/config"
)

func NewGrpcPromObj(cusMetrics []prometheus.Collector) (*prometheus.Registry, *gp.ServerMetrics, error) {
	reg := prometheus.NewRegistry()
	grpcMetrics := gp.NewServerMetrics()
	grpcMetrics.EnableHandlingTimeHistogram()
	cusMetrics = append(cusMetrics, grpcMetrics, collectors.NewGoCollector())
	reg.MustRegister(cusMetrics...)
	return reg, grpcMetrics, nil
}

func GetGrpcCusMetrics(registerName string, share *config.Share) []prometheus.Collector {
	switch registerName {
	case share.RpcRegisterName.Meeting:
		return []prometheus.Collector{MeetingCustomCounter}
	default:
		return nil
	}
}
