package base

import (
	config "openmeeting-server/dto"
	"openmeeting-server/internal/infrastructure/cache"
)

func InitMongo() error {
	err := cache.InitMongoClient(&config.Config.Mongo)
	if err != nil {
		return err
	}
	return nil
}
