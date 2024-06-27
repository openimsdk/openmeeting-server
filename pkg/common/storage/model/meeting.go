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
	TimeZone        string `bson:"time_zone"`
	EndDate         int64  `bson:"end_date"`
	RepeatTimes     int32  `bson:"repeat_times"` // repeat_times means times the meeting repeats
	RepeatType      string `bson:"repeat_type"`  // none, daily, weekly, monthly, custom
	UintType        string `bson:"uint_type"`    // only used when repeat_type is custom
	Interval        int32  `bson:"interval"`     // only used when repeat_type is custom
}
