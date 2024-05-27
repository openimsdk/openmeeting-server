package cmd

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/api/admin"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type AdminApiCmd struct {
	*RootCmd
	ctx       context.Context
	configMap map[string]any
	apiConfig *admin.Config
}

func NewAdminApiCmd() *AdminApiCmd {
	var apiConfig admin.Config
	ret := &AdminApiCmd{apiConfig: &apiConfig}
	ret.configMap = map[string]any{
		OpenMeetingAdminAPICfgFileName: &apiConfig.AdminAPI,
		DiscoveryConfigFilename:        &apiConfig.Discovery,
		ShareFileName:                  &apiConfig.Share,
		MongodbConfigFileName:          &apiConfig.Mongo,
		RedisConfigFileName:            &apiConfig.Redis,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *AdminApiCmd) Exec() error {
	return a.Execute()
}

func (a *AdminApiCmd) runE() error {
	return admin.Start(a.ctx, a.Index(), a.apiConfig)
}
