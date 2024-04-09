package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"openmeeting-server/constant"
	config "openmeeting-server/dto"
	"time"
)

var (
	mongoClient *mongo.Client
)

func InitMongoClient(conf *config.MongoConf) error {
	if mongoClient != nil {
		return errors.New("mongo db init again, please check")
	}

	specialerror.AddReplace(mongo.ErrNoDocuments, errs.ErrRecordNotFound)
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if *conf.Uri != "" {
		uri = *conf.Uri
	} else {
		mongodbHosts := ""
		for i, v := range *conf.Address {
			if i == len(*conf.Address)-1 {
				mongodbHosts += v
			} else {
				mongodbHosts += v + ","
			}
		}
		if conf.Username != nil && conf.Password != nil {
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d",
				*conf.Username, *conf.Password, mongodbHosts,
				*conf.Database, *conf.MaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d",
				mongodbHosts, *conf.Database,
				*conf.MaxPoolSize)
		}
	}
	fmt.Println("mongo:", uri)
	var client *mongo.Client
	var err error
	for i := 0; i <= constant.MongoMaxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		cancel()
		if err == nil {
			mongoClient = client
			return nil
		}
		if cmdErr, ok := err.(mongo.CommandError); ok {
			if cmdErr.Code == 13 || cmdErr.Code == 18 {
				return err
			} else {
				fmt.Printf("Failed to connect to MongoDB: %s\n", err)
			}
		}
	}
	return err
}

func GetMongoClient() (*mongo.Client, error) {
	if mongoClient == nil {
		return nil, errors.New("please init mongodb first")
	}
	return mongoClient, nil
}
