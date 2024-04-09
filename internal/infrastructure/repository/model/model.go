package model

type MeetingInfo struct {
	RoomID      string `bson:"room_id"`
	MeetingName string `bson:"meeting_name"`
}

type SignalModel struct {
	SID string `bson:"sid"`
}

type SignalInvitationModel struct {
	SID string `bson:"sid"`
}
