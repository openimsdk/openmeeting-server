package livekit

import (
	"context"
	lksdk "github.com/livekit/server-sdk-go"
	"github.com/openimsdk/tools/log"
)

type CallbackInterface interface {
	OnRoomParticipantConnected(ctx context.Context)
	OnRoomParticipantDisconnected(ctx context.Context)
	OnRoomDisconnected(ctx context.Context, roomID string, sid string)
	OnMeetingDisconnected(ctx context.Context, roomID string)
	OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string)
}

func NewRoomCallback(ctx context.Context, roomID, sID string,
	roomClient *lksdk.RoomServiceClient, cb CallbackInterface) *RoomCallback {
	return &RoomCallback{
		ctx:        ctx,
		roomID:     roomID,
		sID:        sID,
		roomClient: roomClient,
		cb:         cb,
	}
}

type RoomCallback struct {
	userJoin   bool
	sID        string
	roomID     string
	ctx        context.Context
	roomClient *lksdk.RoomServiceClient
	cb         CallbackInterface
}

func (r *RoomCallback) onParticipantConnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "onParticipantConnected", nil)
	r.cb.OnRoomParticipantConnected(r.ctx)
}

func (r *RoomCallback) onParticipantDisconnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "onParticipantDisconnected", nil)
	r.cb.OnRoomParticipantDisconnected(r.ctx)
}

func (r *RoomCallback) onDisconnected() {
	log.ZWarn(r.ctx, "onDisconnected", nil)
	r.cb.OnRoomDisconnected(r.ctx, r.roomID, r.sID)
}

func (r *RoomCallback) onReconnected() {
	log.ZWarn(r.ctx, "onReconnected", nil)
}

func (r *RoomCallback) onReconnecting() {
	log.ZWarn(r.ctx, "onReconnecting", nil)

}
