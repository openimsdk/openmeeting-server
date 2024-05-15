package config

import (
	_ "embed"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
)

//go:embed version
var Version string

type RedisConf struct {
	ClusterMode    *bool     `mapstructure:"clusterMode"`
	Address        *[]string `mapstructure:"address"`
	Username       *string   `mapstructure:"username"`
	Password       *string   `mapstructure:"password"`
	EnablePipeline *bool     `mapstructure:"enablePipeline"`
}

type Redis struct {
	Address        []string `mapstructure:"address"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	EnablePipeline bool     `mapstructure:"enablePipeline"`
	ClusterMode    bool     `mapstructure:"clusterMode"`
	DB             int      `mapstructure:"db"`
	MaxRetry       int      `mapstructure:"MaxRetry"`
}

type Mongo struct {
	URI         string   `mapstructure:"uri"`
	Address     []string `mapstructure:"address"`
	Database    string   `mapstructure:"database"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	MaxPoolSize int      `mapstructure:"maxPoolSize"`
	MaxRetry    int      `mapstructure:"maxRetry"`
}

func (m *Mongo) Build() *mongoutil.Config {
	return &mongoutil.Config{
		Uri:         m.URI,
		Address:     m.Address,
		Database:    m.Database,
		Username:    m.Username,
		Password:    m.Password,
		MaxPoolSize: m.MaxPoolSize,
		MaxRetry:    m.MaxRetry,
	}
}

func (r *Redis) Build() *redisutil.Config {
	return &redisutil.Config{
		ClusterMode: r.ClusterMode,
		Address:     r.Address,
		Username:    r.Username,
		Password:    r.Password,
		DB:          r.DB,
		MaxRetry:    r.MaxRetry,
	}
}

type RpcRegisterName struct {
	Meeting string `mapstructure:"meeting"`
}

type Share struct {
	Secret          string          `mapstructure:"secret"`
	Env             string          `mapstructure:"env"`
	RpcRegisterName RpcRegisterName `mapstructure:"rpcRegisterName"`
	IMAdminUserID   []string        `mapstructure:"imAdminUserID"`
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

//type Log struct {
//	StorageLocation     *string `mapstructure:"storageLocation"`
//	RotationTime        *uint   `mapstructure:"rotationTime"`
//	RemainRotationCount *uint   `mapstructure:"remainRotationCount"`
//	RemainLogLevel      *int    `mapstructure:"remainLogLevel"`
//	IsStdout            *bool   `mapstructure:"isStdout"`
//	IsJson              *bool   `mapstructure:"isJson"`
//	WithStack           *bool   `mapstructure:"withStack"`
//}

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	WithStack           bool   `mapstructure:"withStack"`
}

type Prometheus struct {
	Enable bool  `mapstructure:"enable"`
	Ports  []int `mapstructure:"ports"`
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
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Config struct {
	RedisConfig   Redis
	MongodbConfig Mongo
	EtcdConfig    EtcdConf
	RpcConfig     Meeting
	Share         Share
}
