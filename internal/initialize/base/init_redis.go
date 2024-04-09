package base

import (
	config "openmeeting-server/dto"
	"openmeeting-server/internal/infrastructure/cache"
)

func InitRedis() error {
	if err := cache.InitRedis(config.Config.Redis); err != nil {
		return err
	}
	return nil
}
