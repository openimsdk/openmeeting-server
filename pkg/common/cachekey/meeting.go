package cachekey

const (
	MeetingInfoKey       = "MEETING_INFO:"
	GenerateMeetingIDKey = "GENERATE_MEETING_ID_KEY"
)

func GetMeetingInfoKey(meetingID string) string {
	return MeetingInfoKey + meetingID
}
