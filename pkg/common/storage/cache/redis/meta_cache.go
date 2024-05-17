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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw/specialerror"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

const (
	scanCount     = 3000
	maxRetryTimes = 5
	retryInterval = time.Millisecond * 100
)

var (
	once      sync.Once
	subscribe map[string][]string
)

var errIndex = errs.New("err index")

func NewMetaCacheRedis(rcClient *rockscache.Client, keys ...string) cache.Meta {
	return &metaCacheRedis{rcClient: rcClient, keys: keys, maxRetryTimes: maxRetryTimes, retryInterval: retryInterval}
}

type metaCacheRedis struct {
	rcClient      *rockscache.Client
	keys          []string
	maxRetryTimes int
	retryInterval time.Duration
	redisClient   redis.UniversalClient
}

func (m *metaCacheRedis) Copy() cache.Meta {
	var keys []string
	if len(m.keys) > 0 {
		keys = make([]string, 0, len(m.keys)*2)
		keys = append(keys, m.keys...)
	}
	return &metaCacheRedis{
		rcClient:      m.rcClient,
		keys:          keys,
		maxRetryTimes: m.maxRetryTimes,
		retryInterval: m.retryInterval,
		redisClient:   m.redisClient,
	}
}

func (m *metaCacheRedis) SetRawRedisClient(cli redis.UniversalClient) {
	m.redisClient = cli
}

func (m *metaCacheRedis) ExecDel(ctx context.Context, distinct ...bool) error {
	if len(distinct) > 0 && distinct[0] {
		m.keys = datautil.Distinct(m.keys)
	}
	if len(m.keys) > 0 {

		for _, key := range m.keys {
			for i := 0; i < m.maxRetryTimes; i++ {
				if err := m.rcClient.TagAsDeleted(key); err != nil {
					log.ZError(ctx, "delete cache failed", err, "key", key)
					time.Sleep(m.retryInterval)
					continue
				}
				break
			}
		}
	}
	return nil
}

func (m *metaCacheRedis) DelKey(ctx context.Context, key string) error {
	return m.rcClient.TagAsDeleted2(ctx, key)
}

func (m *metaCacheRedis) AddKeys(keys ...string) {
	m.keys = append(m.keys, keys...)
}

func (m *metaCacheRedis) ClearKeys() {
	m.keys = []string{}
}

func (m *metaCacheRedis) GetPreDelKeys() []string {
	return m.keys
}

func GetDefaultOpt() rockscache.Options {
	opts := rockscache.NewDefaultOptions()
	opts.StrongConsistency = true
	opts.RandomExpireAdjustment = 0.2

	return opts
}

func getCache[T any](ctx context.Context, rcClient *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rcClient.Fetch2(ctx, key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", errs.WrapMsg(err, "marshal failed")
		}
		write = true

		return string(bs), nil
	})
	if err != nil {
		return t, errs.Wrap(err)
	}
	if write {
		return t, nil
	}
	if v == "" {
		return t, errs.ErrRecordNotFound.WrapMsg("cache is not found")
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		errInfo := fmt.Sprintf("cache json.Unmarshal failed, key:%s, value:%s, expire:%s", key, v, expire)
		return t, errs.WrapMsg(err, errInfo)
	}

	return t, nil
}

func batchGetCache2[T any, K comparable](
	ctx context.Context,
	rcClient *rockscache.Client,
	expire time.Duration,
	keys []K,
	keyFn func(key K) string,
	fns func(ctx context.Context, key K) (T, error),
) ([]T, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	res := make([]T, 0, len(keys))
	for _, key := range keys {
		val, err := getCache(ctx, rcClient, keyFn(key), expire, func(ctx context.Context) (T, error) {
			return fns(ctx, key)
		})
		if err != nil {
			if errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err))) {
				continue
			}
			return nil, errs.Wrap(err)
		}
		res = append(res, val)
	}

	return res, nil
}
