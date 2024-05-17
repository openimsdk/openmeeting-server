package redis

import (
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	meetingExpireTime = time.Second * 60 * 60 * 12
)

type Meeting struct {
	cache.Meta
	rdb        redis.UniversalClient
	meetingDB  database.Meeting
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewMeeting(rdb redis.UniversalClient, meetingDB database.Meeting, options rockscache.Options) cache.Meeting {
	rcClient := rockscache.NewClient(rdb, options)
	mc := NewMetaCacheRedis(rcClient)
	mc.SetRawRedisClient(rdb)
	return &Meeting{
		rdb:        rdb,
		Meta:       NewMetaCacheRedis(rcClient),
		meetingDB:  meetingDB,
		expireTime: meetingExpireTime,
		rcClient:   rcClient,
	}
}

func (m *Meeting) NewCache() cache.Meeting {
	return &Meeting{
		rdb:        m.rdb,
		meetingDB:  m.meetingDB,
		expireTime: m.expireTime,
		rcClient:   m.rcClient,
		Meta:       m.Copy(),
	}
}
