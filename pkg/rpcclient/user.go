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

package rpcclient

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/user"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
	"strings"
)

// User represents a structure holding connection details for the User RPC client.
type User struct {
	conn   grpc.ClientConnInterface
	Client user.UserClient
	Discov discovery.SvcDiscoveryRegistry
}

// NewUser initializes and returns a User instance based on the provided service discovery registry.
func NewUser(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *User {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := user.NewUserClient(conn)
	return &User{Discov: discov, Client: client,
		conn: conn,
	}
}

// UserRpcClient represents the structure for a User RPC client.
type UserRpcClient User

// NewUserRpcClientByUser initializes a UserRpcClient based on the provided User instance.
func NewUserRpcClientByUser(user *User) *UserRpcClient {
	rpc := UserRpcClient(*user)
	return &rpc
}

// NewUserRpcClient initializes a UserRpcClient based on the provided service discovery registry.
func NewUserRpcClient(client discovery.SvcDiscoveryRegistry, rpcRegisterName string,
	imAdminUserID []string) UserRpcClient {
	return UserRpcClient(*NewUser(client, rpcRegisterName))
}

// GetUsersInfo retrieves information for multiple users based on their user IDs.
func (u *UserRpcClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*user.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*user.UserInfo{}, nil
	}
	resp, err := u.Client.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(userIDs, datautil.Slice(resp.UsersInfo, func(e *user.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, servererrs.ErrUserIDNotFound.WrapMsg(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

// GetUserInfo retrieves information for a single user based on the provided user ID.
func (u *UserRpcClient) GetUserInfo(ctx context.Context, userID string) (*user.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

// GetUsersInfoMap retrieves a map of user information indexed by their user IDs.
func (u *UserRpcClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*user.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(users, func(e *user.UserInfo) string {
		return e.UserID
	}), nil
}

// GetPublicUserInfos retrieves public information for multiple users based on their user IDs.
func (u *UserRpcClient) GetPublicUserInfos(
	ctx context.Context,
	userIDs []string,
	complete bool,
) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.Slice(users, func(e *user.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
		}
	}), nil
}

// GetPublicUserInfo retrieves public information for a single user based on the provided user ID.
func (u *UserRpcClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

// GetPublicUserInfoMap retrieves a map of public user information indexed by their user IDs.
func (u *UserRpcClient) GetPublicUserInfoMap(
	ctx context.Context,
	userIDs []string,
	complete bool,
) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}
