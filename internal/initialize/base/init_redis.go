package base

import (
	"openmeeting-server/internal/infrastructure/cache"
	"openmeeting-server/pkg/common/config"
)

func InitRedis() error {
	if err := cache.InitRedis(config.Config.Redis); err != nil {
		return err
	}
	return nil
}
