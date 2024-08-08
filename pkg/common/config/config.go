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

package config

import (
	_ "embed"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
)

//go:embed version
var Version string

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	IsSimplify          bool   `mapstructure:"isSimplify"`
	WithStack           bool   `mapstructure:"withStack"`
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

type Share struct {
	RpcRegisterName RpcRegisterName `mapstructure:"rpcRegisterName"`
}

type API struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
	Api    struct {
		ListenIP string `mapstructure:"listenIP"`
		Ports    []int  `mapstructure:"ports"`
	} `mapstructure:"api"`
	Prometheus struct {
		Enable     bool   `mapstructure:"enable"`
		Ports      []int  `mapstructure:"ports"`
		GrafanaURL string `mapstructure:"grafanaURL"`
	} `mapstructure:"prometheus"`
}

type Prometheus struct {
	Enable bool  `mapstructure:"enable"`
	Ports  []int `mapstructure:"ports"`
}

type User struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Token struct {
		Secret  string `mapstructure:"secret"`
		Expires int    `mapstructure:"expire"`
	} `mapstructure:"token"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Meeting struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type RTC struct {
	URL       []string `mapstructure:"url"`
	ApiKey    string   `mapstructure:"apiKey"`
	ApiSecret string   `mapstructure:"apiSecret"`
	InnerURL  string   `mapstructure:"innerURL"`
}

type Record struct {
	Enable bool `yaml:"enable"`
}

type S3 struct {
	AccessKey string `mapstructure:"access_key"`
	Secret    string `mapstructure:"secret"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
}

type Upload struct {
	Record    Record `mapstructure:"record"`
	Layout    string `mapstructure:"layout"`
	FrameRate int32  `mapstructure:"frame_rate"`
	S3        S3     `mapstructure:"s3"`
}

type Redis struct {
	Address        []string `mapstructure:"address"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	EnablePipeline bool     `mapstructure:"enablePipeline"`
	ClusterMode    bool     `mapstructure:"clusterMode"`
	DB             int      `mapstructure:"storage"`
	MaxRetry       int      `mapstructure:"MaxRetry"`
}

type RpcRegisterName struct {
	User    string `mapstructure:"user"`
	Meeting string `mapstructure:"meeting"`
}

type Discovery struct {
	Enable string `mapstructure:"enable"`
	Etcd   Etcd   `mapstructure:"etcd"`
}

type Etcd struct {
	RootDirectory string   `mapstructure:"rootDirectory"`
	Address       []string `mapstructure:"address"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
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
