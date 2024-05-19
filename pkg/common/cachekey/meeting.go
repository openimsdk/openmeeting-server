package cachekey

const (
	MeetingInfoKey = "MEETING_INFO:"
)

func GetMeetingInfoKey(meetingID string) string {
	return MeetingInfoKey + meetingID
}
