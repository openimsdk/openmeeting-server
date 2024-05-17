package mgo

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMeetingMongo(db *mongo.Database) (database.Meeting, error) {
	coll := db.Collection("meeting")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "meeting_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &MeetingMgo{coll: coll}, nil
}

type MeetingMgo struct {
	coll *mongo.Collection
}

func (u *MeetingMgo) Create(ctx context.Context, meetings []*model.MeetingInfo) error {
	return mongoutil.InsertMany(ctx, u.coll, meetings)
}

func (u *MeetingMgo) Take(ctx context.Context, meetingID string) (user *model.MeetingInfo, err error) {
	return mongoutil.FindOne[*model.MeetingInfo](ctx, u.coll, bson.M{"meeting_id": meetingID})
}
