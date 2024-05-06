package model

import "time"

type MeetingInfo struct {
	MeetingID     string    `bson:"meeting_id"`
	MeetingName   string    `bson:"meeting_name"`
	HostUserID    string    `bson:"host_user_id"`
	StartTime     time.Time `bson:"start_time"`
	EndTime       time.Time `bson:"end_time"`
	Duration      int64     `bson:"duration"`
	Status        int64     `bson:"status"`
	CreatorUserID string    `bson:"creator_user_id"`
	CreateTime    time.Time `bson:"create_time"`
	UpdateTime    time.Time `bson:"update_time"`
}
