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

package user

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/prommetrics"
	mgo2 "github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/tools/db/redisutil"

	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"

	"github.com/openimsdk/openmeeting-server/pkg/common/convert"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	pbuser "github.com/openimsdk/openmeeting-server/pkg/protocol/user"
	"github.com/openimsdk/tools/db/mongoutil"
	registry "github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
	"strings"
)

type userServer struct {
	userStorageHandler controller.User
	RegisterCenter     registry.SvcDiscoveryRegistry
	config             *Config
}

type Config struct {
	Rpc       config.User
	Redis     config.Redis
	Mongo     config.Mongo
	Discovery config.Discovery
	Share     config.Share
}

func Start(ctx context.Context, config *Config, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgoCli, err := mongoutil.NewMongoDB(ctx, config.Mongo.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.Redis.Build())
	if err != nil {
		return err
	}

	userDB, err := mgo2.NewUserMongo(mgoCli.GetDB())
	if err != nil {
		return err
	}
	userCache := redis.NewUser(rdb, userDB, redis.GetDefaultOpt())
	database := controller.NewUser(userDB, userCache, mgoCli.GetTx())
	u := &userServer{
		userStorageHandler: database,
		RegisterCenter:     client,
		config:             config,
	}
	pbuser.RegisterUserServer(server, u)
	return nil
}

func (s *userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (resp *pbuser.GetDesignateUsersResp, err error) {
	resp = &pbuser.GetDesignateUsersResp{}
	users, err := s.userStorageHandler.FindWithError(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}

	resp.UsersInfo = convert.UsersDB2Pb(users)
	return resp, nil
}

func (s *userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if len(req.Users) == 0 {
		return nil, errs.ErrArgs.WrapMsg("users is empty")
	}

	if datautil.DuplicateAny(req.Users, func(e *pbuser.UserInfo) string { return e.UserID }) {
		return nil, errs.ErrArgs.WrapMsg("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, user := range req.Users {
		if user.UserID == "" {
			return nil, errs.ErrArgs.WrapMsg("userID is empty")
		}
		if strings.Contains(user.UserID, ":") {
			return nil, errs.ErrArgs.WrapMsg("userID contains ':' is invalid userID")
		}
		userIDs = append(userIDs, user.UserID)
	}
	users := make([]*model.User, 0, len(req.Users))
	for _, user := range req.Users {
		users = append(users, &model.User{
			UserID:   user.UserID,
			Nickname: user.Nickname,
		})
	}
	if err := s.userStorageHandler.Create(ctx, users); err != nil {
		return nil, err
	}

	prommetrics.UserRegisterCounter.Inc()

	return resp, nil
}
