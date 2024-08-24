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
	OpenMeetingRPCUserCfgFileName    string
	OpenMeetingRPCMeetingCfgFileName string
	RedisConfigFileName              string
	MongodbConfigFileName            string
	DiscoveryConfigFilename          string
	OpenMeetingAPICfgFileName        string
	OpenMeetingAdminAPICfgFileName   string
	LogConfigFileName                string
	ShareFileName                    string
	LiveKitConfigFilename            string
	RecordMeetingFilename            string
)

const envPrefix = "IMENV_"

var ConfigEnvPrefixMap map[string]string

func init() {

	RedisConfigFileName = "redis.yml"
	MongodbConfigFileName = "mongodb.yml"
	OpenMeetingAPICfgFileName = "openmeeting-api.yml"
	OpenMeetingAdminAPICfgFileName = "openmeeting-admin-api.yml"
	OpenMeetingRPCUserCfgFileName = "openmeeting-rpc-user.yml"
	OpenMeetingRPCMeetingCfgFileName = "openmeeting-rpc-meeting.yml"
	DiscoveryConfigFilename = "discovery.yml"
	LogConfigFileName = "log.yml"
	ShareFileName = "share.yml"
	LiveKitConfigFilename = "livekit.yml"
	RecordMeetingFilename = "recorder.yml"

	ConfigEnvPrefixMap = make(map[string]string)
	fileNames := []string{
		RedisConfigFileName,
		MongodbConfigFileName,
		OpenMeetingRPCUserCfgFileName,
		OpenMeetingRPCMeetingCfgFileName,
		DiscoveryConfigFilename,
		OpenMeetingAPICfgFileName,
		OpenMeetingAdminAPICfgFileName,
		LogConfigFileName,
		ShareFileName,
		RecordMeetingFilename,
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
