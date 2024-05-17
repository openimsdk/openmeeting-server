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

package cmd

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/rpc/user"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/prommetrics"
	"github.com/openimsdk/openmeeting-server/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

type UserRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	userConfig *user.Config
}

func NewUserRpcCmd() *UserRpcCmd {
	var userConfig user.Config
	ret := &UserRpcCmd{userConfig: &userConfig}
	ret.configMap = map[string]any{
		OpenIMRPCUserCfgFileName: &userConfig.Rpc,
		RedisConfigFileName:      &userConfig.Redis,
		MongodbConfigFileName:    &userConfig.Mongo,
		ShareFileName:            &userConfig.Share,
		DiscoveryConfigFilename:  &userConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *UserRpcCmd) Exec() error {
	return a.Execute()
}

func (a *UserRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.userConfig.Discovery, &a.userConfig.Rpc.Prometheus, a.userConfig.Rpc.RPC.ListenIP,
		a.userConfig.Rpc.RPC.RegisterIP, a.userConfig.Rpc.RPC.Ports,
		a.Index(), a.userConfig.Share.RpcRegisterName.User, a.userConfig, user.Start, []prometheus.Collector{prommetrics.UserRegisterCounter})
}
