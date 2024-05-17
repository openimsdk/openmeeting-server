package model

// MeetingInfo represents information about a specific meeting.
type MeetingInfo struct {
	MeetingID       string `bson:"meeting_id"`
	Title           string `bson:"title"`
	ScheduledTime   int64  `bson:"scheduled_time"`
	MeetingDuration int64  `bson:"meeting_duration"`
	Password        string `bson:"password"`
	CreatorUserID   string `bson:"creator_user_id"`
	Status          string `bson:"status"`
	StartTime       int64  `bson:"start_time"`
}
