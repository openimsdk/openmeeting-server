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
	"strings"
)

var (
	OpenIMRPCUserCfgFileName string
	RedisConfigFileName      string
	MongodbConfigFileName    string
	DiscoveryConfigFilename  string
	OpenIMAPICfgFileName     string
	LogConfigFileName        string
	ShareFileName            string
)

const envPrefix = "IMENV_"

var ConfigEnvPrefixMap map[string]string

func init() {

	RedisConfigFileName = "redis.yml"
	MongodbConfigFileName = "mongodb.yml"
	OpenIMAPICfgFileName = "openim-api.yml"
	OpenIMRPCUserCfgFileName = "openim-rpc-user.yml"
	DiscoveryConfigFilename = "discovery.yml"
	LogConfigFileName = "log.yml"
	ShareFileName = "share.yml"

	ConfigEnvPrefixMap = make(map[string]string)
	fileNames := []string{
		RedisConfigFileName,
		MongodbConfigFileName,
		OpenIMRPCUserCfgFileName,
		DiscoveryConfigFilename,
		OpenIMAPICfgFileName,
		LogConfigFileName,
		ShareFileName,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = envPrefix + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		ConfigEnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
