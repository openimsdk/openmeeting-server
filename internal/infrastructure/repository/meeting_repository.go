package repository

import (
	"context"
	"github.com/OpenIMSDK/tools/mgoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	config "openmeeting-server/dto"
	"openmeeting-server/internal/infrastructure/cache"
	"openmeeting-server/internal/infrastructure/repository/model"
)

type meetingRepository struct {
	coll *mongo.Collection
}

func NewMeetingRepository() (MeetingInterface, error) {
	client, getClientErr := cache.GetMongoClient()
	if getClientErr != nil {
		return nil, getClientErr
	}
	db := client.Database(*config.Config.Mongo.Database)
	coll := db.Collection("meeting")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "room_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, err
	}
	return &meetingRepository{coll: coll}, nil
}

type MeetingInterface interface {
	DeleteMeetingInfos(ctx context.Context, roomIDs []string) error
	CreateMeetingInfo(ctx context.Context, meetingInfo *model.MeetingInfo) error
	UpdateMeetingInfo(ctx context.Context, roomID string, updateData map[string]any) error
}

func (m *meetingRepository) DeleteMeetingInfos(ctx context.Context, roomIDs []string) error {
	return mgoutil.DeleteMany(ctx, m.coll, bson.M{"room_id": bson.M{"$in": roomIDs}})
}

func (m *meetingRepository) CreateMeetingInfo(ctx context.Context, meetingInfo *model.MeetingInfo) error {
	return mgoutil.InsertMany(ctx, m.coll, []*model.MeetingInfo{meetingInfo})
}

func (m *meetingRepository) UpdateMeetingInfo(ctx context.Context, roomID string, updateData map[string]any) error {
	if len(updateData) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, m.coll, bson.M{"room_id": roomID}, bson.M{"$set": updateData}, false)
}
