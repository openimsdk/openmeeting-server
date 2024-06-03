package rtc

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
)

type MeetingRtc interface {
	GetJoinToken(ctx context.Context, roomID, identity string, metadata *meeting.ParticipantMetaData) (string, string, error)
	CreateRoom(ctx context.Context, roomID, identify string, metaData *meeting.ParticipantMetaData) (sID, token, liveUrl string, err error)
	GetRoomData(ctx context.Context, roomID string) (*meeting.MeetingMetadata, error)
	RoomIsExist(ctx context.Context, roomID string) (string, error)
	UpdateMetaData(ctx context.Context, info *meeting.MeetingMetadata) error
	CloseRoom(ctx context.Context, roomID string) error
	RemoveParticipant(ctx context.Context, roomID, userID string) error
	ToggleMimeStream(ctx context.Context, roomID, userID, mineType string, mute bool) error
}
