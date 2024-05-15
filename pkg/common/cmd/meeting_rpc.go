package cmd

import (
	"context"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
	"openmeeting-server/internal/rpc/meeting"
	startrpc "openmeeting-server/pkg/common"
	"openmeeting-server/pkg/common/config"
)

type MeetingRpcCmd struct {
	*RootCmd
	ctx           context.Context
	configMap     map[string]any
	meetingConfig *config.Config
}

func NewMeetingRpcCmd() *MeetingRpcCmd {
	var meetingConfig config.Config
	ret := &MeetingRpcCmd{meetingConfig: &meetingConfig}
	ret.configMap = map[string]any{
		OpenMeetingRPCCfgFileName: &meetingConfig.RpcConfig,
		RedisConfigFileName:       &meetingConfig.RedisConfig,
		EtcdConfigFileName:        &meetingConfig.EtcdConfig,
		MongodbConfigFileName:     &meetingConfig.MongodbConfig,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.Background()
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *MeetingRpcCmd) Exec() error {
	return a.Execute()
}

func (a *MeetingRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.meetingConfig.RpcConfig.Prometheus, a.meetingConfig.RpcConfig.RPC.ListenIP,
		a.meetingConfig.RpcConfig.RPC.RegisterIP, a.meetingConfig.RpcConfig.RPC.Port, a.Index(),
		a.meetingConfig.Share.RpcRegisterName.Meeting, *a.meetingConfig, &a.meetingConfig.Share, meeting.Start)
}
