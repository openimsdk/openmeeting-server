package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signalInvitationRepository struct {
	coll *mongo.Collection
}

type SignalInvitationInterface interface {
	Delete(ctx context.Context, sids []string) error
}

func NewSignalInvitation(db *mongo.Database) (SignalInvitationInterface, error) {
	coll := db.Collection("signal_invitation")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "sid", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, err
	}
	return &signalInvitationRepository{coll: coll}, nil
}

func (signal *signalInvitationRepository) Delete(ctx context.Context, sids []string) error {
	return nil
}
