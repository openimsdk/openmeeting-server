package rtc

import (
	"context"
	lksdk "github.com/livekit/server-sdk-go"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/tools/log"
)

type CallbackInterface interface {
	OnRoomParticipantConnected(ctx context.Context, userID string)
	OnRoomParticipantDisconnected(ctx context.Context, userID string)
	OnRoomDisconnected(ctx context.Context)
	OnMeetingDisconnected(ctx context.Context, roomID string)
	OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string)
}

func NewRoomCallback(ctx context.Context, roomID, sID string,
	cb CallbackInterface, userRpc *rpcclient.User) *RoomCallback {
	return &RoomCallback{
		ctx:     ctx,
		roomID:  roomID,
		sID:     sID,
		cb:      cb,
		userRpc: userRpc,
	}
}

type RoomCallback struct {
	userJoin bool
	sID      string
	roomID   string
	ctx      context.Context
	cb       CallbackInterface
	userRpc  *rpcclient.User
}

func (r *RoomCallback) OnParticipantConnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "OnParticipantConnected", nil)
	r.cb.OnRoomParticipantConnected(r.ctx, rp.Identity())
}

func (r *RoomCallback) OnParticipantDisconnected(rp *lksdk.RemoteParticipant) {
	log.ZWarn(r.ctx, "OnParticipantDisconnected", nil, "kind:", rp.Kind().String())
	r.cb.OnRoomParticipantDisconnected(r.ctx, rp.Identity())
	// clear user token
	//if _, err := r.userRpc.Client.ClearUserToken(r.ctx, &pbuser.ClearUserTokenReq{UserID: rp.Identity()}); err != nil {
	//	log.ZWarn(r.ctx, "clear user token failed", err, "userID", rp.Identity())
	//}
}

func (r *RoomCallback) OnDisconnected() {
	log.ZWarn(r.ctx, "OnDisconnected", nil)
	r.cb.OnRoomDisconnected(r.ctx)
}

func (r *RoomCallback) OnReconnected() {
	log.ZWarn(r.ctx, "OnReconnected", nil)
}

func (r *RoomCallback) OnReconnecting() {
	log.ZWarn(r.ctx, "OnReconnecting", nil)
}
