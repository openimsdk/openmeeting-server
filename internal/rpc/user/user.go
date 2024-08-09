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
	"errors"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/constant"
	"github.com/openimsdk/openmeeting-server/pkg/common/convert"
	"github.com/openimsdk/openmeeting-server/pkg/common/prommetrics"
	"github.com/openimsdk/openmeeting-server/pkg/common/securetools"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	pbuser "github.com/openimsdk/protocol/openmeeting/user"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
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
	tokenVerify        *token.Token
	meetingRpc         *rpcclient.Meeting
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

	userDB, err := mgo.NewUserMongo(mgoCli.GetDB())
	if err != nil {
		return err
	}
	userCache := redis.NewUser(rdb, userDB, redis.GetDefaultOpt())
	database := controller.NewUser(userDB, userCache, mgoCli.GetTx())
	tokenVerify := token.New(config.Rpc.Token.Expires, config.Rpc.Token.Secret)
	// init rpc client here
	meetingRpc := rpcclient.NewMeeting(client, config.Share.RpcRegisterName.Meeting)

	u := &userServer{
		userStorageHandler: database,
		RegisterCenter:     client,
		config:             config,
		tokenVerify:        tokenVerify,
		meetingRpc:         meetingRpc,
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
		return nil, servererrs.ErrRegisteredAlready.WrapMsg("userID repeated")
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

func (s *userServer) UserLogin(ctx context.Context, req *pbuser.UserLoginReq) (*pbuser.UserLoginResp, error) {
	resp := &pbuser.UserLoginResp{}
	user, err := s.userStorageHandler.GetByAccount(ctx, req.Account)
	if err != nil {
		return resp, servererrs.ErrUserPasswordError.WrapMsg("wrong password or user account")
	}
	saltPasswd := securetools.VerifyPassword(req.Password, user.SaltValue)
	if saltPasswd != user.Password {
		return resp, servererrs.ErrUserPasswordError.WrapMsg("wrong password or user account")
	}
	userToken, err := s.tokenVerify.CreateToken(user.UserID)
	if err != nil {
		return resp, err
	}
	if err := s.userStorageHandler.StoreToken(ctx, user.UserID, userToken); err != nil {
		return resp, err
	}
	resp.UserID = user.UserID
	resp.Token = userToken
	resp.Nickname = user.Nickname
	cleanMsg := &pbmeeting.CleanPreviousMeetingsReq{
		UserID:     user.UserID,
		Reason:     constant.KickOffDuplicatedLogin,
		ReasonCode: int32(pbmeeting.KickOffReason_DuplicatedLogin),
	}

	if _, err := s.meetingRpc.Client.CleanPreviousMeetings(ctx, cleanMsg); err != nil {
		if err != nil {
			return nil, errs.WrapMsg(err, "clean meeting failed")
		}
	}

	return resp, nil
}

func (s *userServer) GetUserToken(ctx context.Context, req *pbuser.GetUserTokenReq) (*pbuser.GetUserTokenResp, error) {
	resp := &pbuser.GetUserTokenResp{}
	userToken, err := s.userStorageHandler.GetToken(ctx, req.UserID)
	if err != nil {
		return resp, servererrs.ErrUserTokenNotFoundErr.WrapMsg("get user token failed")
	}
	resp.Token = userToken
	return resp, nil
}

func (s *userServer) GetUserInfo(ctx context.Context, req *pbuser.GetUserInfoReq) (*pbuser.GetUserInfoResp, error) {
	resp := &pbuser.GetUserInfoResp{}
	userInfo, err := s.userStorageHandler.FindWithError(ctx, []string{req.UserID})
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return resp, servererrs.ErrUserAccountNotFoundErr.WrapMsg("not found user")
		}
		return resp, servererrs.ErrDatabase.WrapMsg("get user failed")
	}
	resp.Account = userInfo[0].Account
	resp.Nickname = userInfo[0].Nickname
	resp.UserID = userInfo[0].UserID
	return resp, nil
}

func (s *userServer) UpdateUserPassword(context.Context, *pbuser.UpdateUserPasswordReq) (*pbuser.UpdateUserPasswordResp, error) {
	resp := &pbuser.UpdateUserPasswordResp{}
	return resp, nil
}

func (s *userServer) ClearUserToken(ctx context.Context, req *pbuser.ClearUserTokenReq) (*pbuser.ClearUserTokenResp, error) {
	resp := &pbuser.ClearUserTokenResp{}

	if err := s.userStorageHandler.StoreToken(ctx, req.UserID, constant.KickOffMeetingMsg); err != nil {
		return resp, errs.WrapMsg(err, "clear user token failed", "user", req.UserID)
	}

	return resp, nil
}

func (s *userServer) UserLogout(ctx context.Context, req *pbuser.LogoutReq) (*pbuser.LogoutResp, error) {
	resp := &pbuser.LogoutResp{}

	if err := s.userStorageHandler.ClearUserToken(ctx, req.UserID); err != nil {
		return resp, errs.WrapMsg(err, "clear token failed")
	}

	// clean previous meeting
	cleanMsg := &pbmeeting.CleanPreviousMeetingsReq{
		UserID:     req.UserID,
		Reason:     constant.KickOffLogout,
		ReasonCode: int32(pbmeeting.KickOffReason_Logout),
	}
	if _, err := s.meetingRpc.Client.CleanPreviousMeetings(ctx, cleanMsg); err != nil {
		if err != nil {
			return nil, errs.WrapMsg(err, "clean meeting failed")
		}
	}

	return resp, nil
}

func (s *userServer) ParseToken(ctx context.Context, req *pbuser.ParseTokenReq) (*pbuser.ParseTokenResp, error) {
	resp := &pbuser.ParseTokenResp{}
	userID, err := s.tokenVerify.GetToken(req.Token)
	if err != nil {
		return resp, err
	}
	userToken, err := s.userStorageHandler.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userToken == constant.KickOffMeetingMsg {
		return nil, servererrs.ErrKickOffMeeting.WrapMsg("kick off meeting, please login again")
	}
	if req.Token != userToken {
		return nil, servererrs.ErrKickOffMeeting.WrapMsg("kick off meeting for login duplicated, please login again")
	}
	return &pbuser.ParseTokenResp{
		UserID: userID,
	}, nil
}
