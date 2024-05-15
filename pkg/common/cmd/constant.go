package cmd

var (
	FileName                  string
	OpenMeetingRPCCfgFileName string
	LocalCacheConfigFileName  string
	RedisConfigFileName       string
	EtcdConfigFileName        string
	MongodbConfigFileName     string
	MinioConfigFileName       string
	LogConfigFileName         string
)

var ConfigEnvPrefixMap map[string]string

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
