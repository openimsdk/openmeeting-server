package rtc

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
	roomClient *lksdk.RoomServiceClient) *RoomCallback {
	return &RoomCallback{
		ctx:        ctx,
		roomID:     roomID,
		sID:        sID,
		roomClient: roomClient,
	}
}

type RoomCallback struct {
	userJoin   bool
	sID        string
	roomID     string
	ctx        context.Context
	roomClient *lksdk.RoomServiceClient
}

func (r *RoomCallback) OnParticipantConnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "OnParticipantConnected", nil)
}

func (r *RoomCallback) OnParticipantDisconnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "OnParticipantDisconnected", nil)
}

func (r *RoomCallback) OnDisconnected() {
	log.ZWarn(r.ctx, "OnDisconnected", nil)
}

func (r *RoomCallback) OnReconnected() {
	log.ZWarn(r.ctx, "OnReconnected", nil)
}

func (r *RoomCallback) OnReconnecting() {
	log.ZWarn(r.ctx, "OnReconnecting", nil)
}
