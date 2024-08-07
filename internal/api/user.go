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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/openmeeting/user"
	"github.com/openimsdk/tools/a2r"
)

type UserApi struct {
	Client user.UserClient
}

func NewUserApi(client user.UserClient) *UserApi {
	return &UserApi{Client: client}
}

func (u *UserApi) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.Client, c)
}

func (u *UserApi) UserLogin(c *gin.Context) {
	a2r.Call(user.UserClient.UserLogin, u.Client, c)
}

func (u *UserApi) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.Client, c)
}

func (u *UserApi) UpdateUserPassword(c *gin.Context) {
}

func (u *UserApi) UserLogout(c *gin.Context) {
	a2r.Call(user.UserClient.UserLogout, u.Client, c)
}
