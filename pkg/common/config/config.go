package config

import config1 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

type RedisConf struct {
	ClusterMode    *bool     `mapstructure:"clusterMode"`
	Address        *[]string `mapstructure:"address"`
	Username       *string   `mapstructure:"username"`
	Password       *string   `mapstructure:"password"`
	EnablePipeline *bool     `mapstructure:"enablePipeline"`
}

type MongoConf struct {
	Uri         *string   `mapstructure:"uri"`
	Address     *[]string `mapstructure:"address"`
	Database    *string   `mapstructure:"database"`
	Username    *string   `mapstructure:"username"`
	Password    *string   `mapstructure:"password"`
	MaxPoolSize *int      `mapstructure:"maxPoolSize"`
}

type EtcdConf struct {
	Address  *[]string `mapstructure:"address"`
	Ttl      *int      `mapstructure:"ttl"`
	Lease    *int      `mapstructure:"lease"`
	Username *string   `mapstructure:"username"`
	Password *string   `mapstructure:"password"`
	Timeout  *int      `mapstructure:"timeout"`
}

type Log struct {
	StorageLocation     *string `mapstructure:"storageLocation"`
	RotationTime        *uint   `mapstructure:"rotationTime"`
	RemainRotationCount *uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      *int    `mapstructure:"remainLogLevel"`
	IsStdout            *bool   `mapstructure:"isStdout"`
	IsJson              *bool   `mapstructure:"isJson"`
	WithStack           *bool   `mapstructure:"withStack"`
}

type Meeting struct {
	RPC struct {
		ListenIP   string `mapstructure:"listenIP"`
		RegisterIP string `mapstructure:"registerIP"`
		Name       string `mapstructure:"name"`
		Port       []int  `mapstructure:"port"`
	} `mapstructure:"rpc"`
	//`mapstructure:"log"`
	RTC struct {
		URL       []string `mapstructure:"url"`
		ApiKey    string   `mapstructure:"apiKey"`
		ApiSecret string   `mapstructure:"apiSecret"`
		InnerURL  string   `mapstructure:"innerURL"`
	} `mapstructure:"rtc"`
	Prometheus config1.Prometheus `mapstructure:"prometheus"`
}

type Config struct {
	RedisConfig   config1.Redis
	MongodbConfig config1.Mongo
	EtcdConfig    EtcdConf
	RpcConfig     Meeting
	Share         config1.Share
}
