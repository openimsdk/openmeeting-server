package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signalRepository struct {
	coll *mongo.Collection
}

func NewSignal(db *mongo.Database) (SignalInterface, error) {
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
