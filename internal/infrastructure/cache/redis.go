package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw/specialerror"
	"openmeeting-server/constant"
	"openmeeting-server/pkg/common/config"

	"github.com/redis/go-redis/v9"
	"time"
)

var (
	// singleton pattern.
	redisClient *redis.UniversalClient
)

func InitRedis(conf config.RedisConf) error {
	if redisClient != nil {
		return errors.New("redis is already init, please check")
	}

	if len(*conf.Address) == 0 {
		return errors.New("redis address is empty")
	}
	specialerror.AddReplace(redis.Nil, errs.ErrRecordNotFound)
	var rdb redis.UniversalClient
	if len(*conf.Address) > 1 || conf.ClusterMode != nil {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      *conf.Address,
			Username:   *conf.Username,
			Password:   *conf.Password, // no password set
			PoolSize:   50,
			MaxRetries: constant.RedisMaxRetry,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:       (*conf.Address)[0],
			Username:   *conf.Username,
			Password:   *conf.Password, // no password set
			DB:         0,              // use default DB
			PoolSize:   100,            // connection pool size
			MaxRetries: constant.RedisMaxRetry,
		})
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping %w", err)
	}

	redisClient = &rdb
	return err
}

func GetRedisClient() (*redis.UniversalClient, error) {
	if redisClient == nil {
		return nil, errors.New("redis has not init")
	}
	return redisClient, nil
}
