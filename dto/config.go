package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type RedisConf struct {
	ClusterMode    *bool     `yaml:"clusterMode"`
	Address        *[]string `yaml:"address"`
	Username       *string   `yaml:"username"`
	Password       *string   `yaml:"password"`
	EnablePipeline *bool     `yaml:"enablePipeline"`
}

type MongoConf struct {
	Uri         *string   `yaml:"uri"`
	Address     *[]string `yaml:"address"`
	Database    *string   `yaml:"database"`
	Username    *string   `yaml:"username"`
	Password    *string   `yaml:"password"`
	MaxPoolSize *int      `yaml:"maxPoolSize"`
}

type EtcdConf struct {
	Address  *[]string `yaml:"address"`
	Ttl      *int      `yaml:"ttl"`
	Lease    *int      `yaml:"lease"`
	Username *string   `yaml:"username"`
	Password *string   `yaml:"password"`
	Timeout  *int      `yaml:"timeout"`
}

var Config struct {
	Mongo MongoConf `yaml:"mongo"`
	Redis RedisConf `yaml:"redis"`
	Etcd  EtcdConf  `yaml:"etcd"`

	RPC struct {
		ListenIP   string `yaml:"listenIP"`
		RegisterIP string `yaml:"registerIP"`
		Name       string `yaml:"name"`
		Port       []int  `yaml:"port"`
	}

	RTC struct {
		URL       []string `yaml:"url"`
		ApiKey    string   `yaml:"apiKey"`
		ApiSecret string   `yaml:"apiSecret"`
		InnerURL  string   `yaml:"innerURL"`
	} `yaml:"rtc"`

	Log struct {
		StorageLocation     *string `yaml:"storageLocation"`
		RotationTime        *uint   `yaml:"rotationTime"`
		RemainRotationCount *uint   `yaml:"remainRotationCount"`
		RemainLogLevel      *int    `yaml:"remainLogLevel"`
		IsStdout            *bool   `yaml:"isStdout"`
		IsJson              *bool   `yaml:"isJson"`
		WithStack           *bool   `yaml:"withStack"`
	} `yaml:"log"`
}

func Parse(fp string) error {
	data, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}
	return nil
}
