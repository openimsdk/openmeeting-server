package livekit

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting/rtc"
	"github.com/openimsdk/tools/log"
)

type CallbackLiveKit struct {
	roomID  string
	liveKit *LiveKit
}

func NewRTC(roomID string, liveKit *LiveKit) rtc.CallbackInterface {
	return &CallbackLiveKit{
		roomID:  roomID,
		liveKit: liveKit,
	}
}

func (r *CallbackLiveKit) OnRoomParticipantConnected(ctx context.Context, userID string) {
	log.ZDebug(ctx, "OnRoomParticipantConnected", nil)
}

func (r *CallbackLiveKit) OnRoomParticipantDisconnected(ctx context.Context, userID string) {
	log.ZWarn(ctx, "OnRoomParticipantDisconnected", nil)
	if err := r.liveKit.RemoveParticipant(ctx, r.roomID, userID); err != nil {
		log.ZWarn(ctx, "remove participant failed", err)
	}
}

func (r *CallbackLiveKit) OnRoomDisconnected(ctx context.Context) {
	log.ZWarn(ctx, "OnRoomDisconnected", nil, r.roomID)
	participants, err := r.liveKit.ListParticipants(ctx, r.roomID)
	if err != nil {
		log.ZWarn(ctx, "remove participant failed", err)
		return
	}
	for _, p := range participants {
		if err := r.liveKit.RemoveParticipant(ctx, r.roomID, p.Identity); err != nil {
			log.ZWarn(ctx, "remove participant failed", err)
		}
	}
}

func (r *CallbackLiveKit) OnMeetingDisconnected(ctx context.Context, roomID string) {
	log.ZWarn(ctx, "OnMeetingDisconnected", nil, roomID)
}

func (r *CallbackLiveKit) OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string) {
	log.ZWarn(ctx, "OnMeetingUnmute", nil, roomID)
}
