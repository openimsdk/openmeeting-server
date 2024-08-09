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

package redis

import (
	"context"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/openmeeting-server/pkg/common/cachekey"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache"
	"github.com/redis/go-redis/v9"
)

const (
	userExpireTime            = time.Second * 60 * 60 * 12
	olineStatusKey            = "ONLINE_STATUS:"
	userOlineStatusExpireTime = time.Second * 60 * 60 * 24
	statusMod                 = 501
)

type User struct {
	cache.Meta
	rdb        redis.UniversalClient
	userDB     database.User
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewUser(rdb redis.UniversalClient, userDB database.User, options rockscache.Options) cache.User {
	rcClient := rockscache.NewClient(rdb, options)
	mc := NewMetaCacheRedis(rcClient)
	//u := localCache.User
	//mc.SetTopic(u.Topic)
	mc.SetRawRedisClient(rdb)
	return &User{
		rdb:        rdb,
		Meta:       NewMetaCacheRedis(rcClient),
		userDB:     userDB,
		expireTime: userExpireTime,
		rcClient:   rcClient,
	}
}

func (u *User) NewCache() cache.User {
	return &User{
		rdb:        u.rdb,
		userDB:     u.userDB,
		expireTime: u.expireTime,
		rcClient:   u.rcClient,
		Meta:       u.Copy(),
	}
}

func (u *User) getUserInfoKey(userID string) string {
	return cachekey.GetUserInfoKey(userID)
}

func (u *User) getUserGlobalRecvMsgOptKey(userID string) string {
	return cachekey.GetUserGlobalRecvMsgOptKey(userID)
}

func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model.User, error) {
	return batchGetCache2(ctx, u.rcClient, u.expireTime, userIDs, func(userID string) string {
		return u.getUserInfoKey(userID)
	}, func(ctx context.Context, userID string) (*model.User, error) {
		return u.userDB.Take(ctx, userID)
	})
}

func (u *User) DelUsersInfo(userIDs ...string) cache.User {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, u.getUserInfoKey(userID))
	}
	cache := u.NewCache()
	cache.AddKeys(keys...)

	return cache
}

func (u *User) GetUserByAccount(ctx context.Context, account string) (*model.User, error) {
	return getCache(ctx, u.rcClient, u.getUserInfoKey(account), u.expireTime, func(ctx context.Context) (*model.User, error) {
		return u.userDB.TakeByAccount(ctx, account)
	})
}

func (u *User) CacheUserToken(ctx context.Context, userID, userToken string) error {
	return errs.Wrap(u.rdb.Set(ctx, cachekey.GetUserTokenKey(userID), userToken, u.expireTime).Err())
}

func (u *User) GetUserToken(ctx context.Context, userID string) (string, error) {
	token, err := u.rdb.Get(ctx, cachekey.GetUserTokenKey(userID)).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return token, nil
}

func (u *User) ClearUserToken(ctx context.Context, userID string) error {
	return errs.Wrap(u.rdb.Del(ctx, cachekey.GetUserTokenKey(userID)).Err())
}

func (u *User) GenerateUserID(ctx context.Context) (string, error) {
	index, err := u.rdb.Incr(ctx, cachekey.GetGenerateUserIDKey()).Result()
	if err != nil {
		return "", errs.WrapMsg(err, "incr key failed from redis")
	}
	return fmt.Sprintf("%08d", index), nil
}

type Comparable interface {
	~int | ~string | ~float64 | ~int32
}
