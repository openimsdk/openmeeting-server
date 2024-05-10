package base

import (
	"openmeeting-server/internal/infrastructure/cache"
	"openmeeting-server/pkg/common/config"
)

func InitMongo() error {
	err := cache.InitMongoClient(&config.Config.Mongo)
	if err != nil {
		return err
	}
	return nil
}
