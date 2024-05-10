package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"openmeeting-server/internal/infrastructure/cache"
	"openmeeting-server/pkg/common/config"
)

type signalRepository struct {
	coll *mongo.Collection
}

func NewSignal() (SignalInterface, error) {
	client, getClientErr := cache.GetMongoClient()
	if getClientErr != nil {
		return nil, getClientErr
	}
	db := client.Database(*config.Config.Mongo.Database)
	coll := db.Collection("signal")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "sid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, err
	}
	return &signalRepository{coll: coll}, nil
}

type SignalInterface interface {
	Delete(ctx context.Context, sids []string) error
}

func (signal *signalRepository) Delete(ctx context.Context, sids []string) error {
	return nil
}
