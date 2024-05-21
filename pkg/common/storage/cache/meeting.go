package cache

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
)

type Meeting interface {
	Meta
	NewCache() Meeting
	GetMeetingByID(ctx context.Context, meetingID string) (*model.MeetingInfo, error)
	DelMeeting(meetingIDs ...string) Meeting
}
