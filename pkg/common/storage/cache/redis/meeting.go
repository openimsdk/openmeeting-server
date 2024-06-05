package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/openmeeting-server/pkg/common/cachekey"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
	"strconv"
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

func (m *Meeting) getMeetingInfoKey(meetingID string) string {
	return cachekey.GetUserInfoKey(meetingID)
}

func (m *Meeting) GetMeetingByID(ctx context.Context, meetingID string) (*model.MeetingInfo, error) {
	return getCache(ctx, m.rcClient, m.getMeetingInfoKey(meetingID), m.expireTime, func(ctx context.Context) (*model.MeetingInfo, error) {
		return m.meetingDB.Take(ctx, meetingID)
	})
}

func (m *Meeting) DelMeeting(meetingIDs ...string) cache.Meeting {
	keys := make([]string, 0, len(meetingIDs))
	for _, meetingID := range meetingIDs {
		keys = append(keys, m.getMeetingInfoKey(meetingID))
	}
	newMeetingCache := m.NewCache()
	newMeetingCache.AddKeys(keys...)

	return newMeetingCache
}

func (m *Meeting) GenerateMeetingID(ctx context.Context) (string, error) {
	value, err := m.rdb.Get(ctx, cachekey.GenerateMeetingIDKey).Result()
	var meetingID int
	if err != nil {
		if errors.Is(err, redis.Nil) {
			value = "1"
		} else {
			return "", errs.WrapMsg(err, "get key failed", cachekey.GenerateMeetingIDKey)
		}
	}
	meetingID, err = strconv.Atoi(value)
	if err != nil {
		return "", errs.WrapMsg(err, "string to integer failed")
	}
	meetingID += 1
	strMeetingID := fmt.Sprintf("%09d", meetingID)
	_ = m.rdb.Set(ctx, cachekey.GenerateMeetingIDKey, strMeetingID, 0)
	return strMeetingID, nil
}
