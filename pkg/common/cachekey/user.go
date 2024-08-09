// Copyright © 2024 OpenIM. All rights reserved.
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

package cachekey

const (
	UserInfoKey             = "USER_INFO:"
	UserTokenKey            = "USER_TOKEN:"
	UserGlobalRecvMsgOptKey = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
	GenerateUserIDKey       = "GENERATE_USER_ID_KEY"
)

func GetUserInfoKey(userID string) string {
	return UserInfoKey + userID
}

func GetUserTokenKey(userID string) string {
	return UserTokenKey + userID
}

func GetGenerateUserIDKey() string {
	return GenerateUserIDKey
}

func GetUserGlobalRecvMsgOptKey(userID string) string {
	return UserGlobalRecvMsgOptKey + userID
}
