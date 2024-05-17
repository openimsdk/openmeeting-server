package cmd

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/prommetrics"
	"github.com/openimsdk/openmeeting-server/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

type MeetingRpcCmd struct {
	*RootCmd
	ctx           context.Context
	configMap     map[string]any
	meetingConfig *meeting.Config
}

func NewMeetingRpcCmd() *MeetingRpcCmd {
	var meetingConfig meeting.Config
	ret := &MeetingRpcCmd{meetingConfig: &meetingConfig}
	ret.configMap = map[string]any{
		OpenIMRPCUserCfgFileName: &meetingConfig.Rpc,
		RedisConfigFileName:      &meetingConfig.Redis,
		MongodbConfigFileName:    &meetingConfig.Mongo,
		ShareFileName:            &meetingConfig.Share,
		DiscoveryConfigFilename:  &meetingConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *MeetingRpcCmd) Exec() error {
	return a.Execute()
}

func (a *MeetingRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.meetingConfig.Discovery, &a.meetingConfig.Rpc.Prometheus, a.meetingConfig.Rpc.RPC.ListenIP,
		a.meetingConfig.Rpc.RPC.RegisterIP, a.meetingConfig.Rpc.RPC.Ports,
		a.Index(), a.meetingConfig.Share.RpcRegisterName.User, a.meetingConfig, meeting.Start, []prometheus.Collector{prommetrics.UserRegisterCounter})
}
