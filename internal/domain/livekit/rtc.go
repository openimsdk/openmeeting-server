package livekit

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
)

type RTC struct {
	ctx context.Context
}

func NewRTC(ctx context.Context) CallbackInterface {
	return &RTC{
		ctx: ctx,
	}
}

func (r *RTC) OnRoomParticipantConnected(ctx context.Context) {
	log.ZWarn(r.ctx, "OnRoomParticipantConnected", nil)
}

func (r *RTC) OnRoomParticipantDisconnected(ctx context.Context) {
	log.ZWarn(r.ctx, "OnRoomParticipantDisconnected", nil)
}

func (r *RTC) OnRoomDisconnected(ctx context.Context, roomID string, sid string) {
	log.ZWarn(r.ctx, "OnRoomDisconnected", nil, roomID, sid)
}

func (r *RTC) OnMeetingDisconnected(ctx context.Context, roomID string) {
	log.ZWarn(r.ctx, "OnMeetingDisconnected", nil, roomID)

}

func (r *RTC) OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string) {

	log.ZWarn(r.ctx, "OnMeetingUnmute", nil, roomID)

}
