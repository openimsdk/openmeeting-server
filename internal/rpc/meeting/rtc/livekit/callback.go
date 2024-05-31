package livekit

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting/rtc"
	"github.com/openimsdk/tools/log"
)

type CallbackLiveKit struct {
}

func NewRTC(ctx context.Context) rtc.CallbackInterface {
	return &CallbackLiveKit{}
}

func (r *CallbackLiveKit) OnRoomParticipantConnected(ctx context.Context) {
	log.ZDebug(ctx, "OnRoomParticipantConnected", nil)
}

func (r *CallbackLiveKit) OnRoomParticipantDisconnected(ctx context.Context) {
	log.ZWarn(ctx, "OnRoomParticipantDisconnected", nil)
}

func (r *CallbackLiveKit) OnRoomDisconnected(ctx context.Context, roomID string, sid string) {
	log.ZWarn(ctx, "OnRoomDisconnected", nil, roomID, sid)
}

func (r *CallbackLiveKit) OnMeetingDisconnected(ctx context.Context, roomID string) {
	log.ZWarn(ctx, "OnMeetingDisconnected", nil, roomID)
}

func (r *CallbackLiveKit) OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string) {
	log.ZWarn(ctx, "OnMeetingUnmute", nil, roomID)
}
