package livekit

import (
	"context"
	"encoding/json"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting/rtc"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/twitchtv/twirp"
	"sync/atomic"
	"time"
)

type LiveKit struct {
	roomClient *lksdk.RoomServiceClient
	index      uint64
	conf       *config.RTC
}

func NewLiveKit(conf *config.RTC) rtc.MeetingRtc {
	return &LiveKit{
		index:      0,
		conf:       conf,
		roomClient: lksdk.NewRoomServiceClient(conf.InnerURL, conf.ApiKey, conf.ApiSecret),
	}
}

func (x *LiveKit) GetJoinToken(ctx context.Context, roomID, identity string) (string, string, error) {
	log.ZDebug(ctx, "getJoinToken", "roomID", roomID, "identity", identity)
	canPublish := true
	canSubscribe := true
	canPublishData := true
	// 配置里面的
	at := auth.NewAccessToken(x.conf.ApiKey, x.conf.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:       true,
		Room:           roomID,
		CanPublish:     &canPublish,
		CanSubscribe:   &canSubscribe,
		CanPublishData: &canPublishData,
	}
	// 生成邀请者房间的jwt
	at.AddGrant(grant).
		SetIdentity(identity).
		SetName("participant-name").
		SetValidFor(time.Hour)
	jwt, err := at.ToJWT()
	if err != nil {
		return "", "", errs.WrapMsg(err, "at.ToJWT failed")
	}
	log.ZDebug(ctx, "getJoinToken", "jwt", jwt)
	return jwt, x.getLiveURL(), nil
}

func (x *LiveKit) CreateRoom(ctx context.Context, roomID string) (sID, token, liveUrl string, err error) {
	room, err := x.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomID,
		EmptyTimeout:    3,
		MaxParticipants: 10000,
	})
	if err != nil {
		log.ZError(ctx, "Marshal failed", err)
		return "", "", "", errs.Wrap(err)
	}
	callback := rtc.NewRoomCallback(mcontext.NewCtx("room_callback_"+mcontext.GetOperationID(ctx)), roomID, room.Sid, x.roomClient)
	roomCallback := &lksdk.RoomCallback{
		ParticipantCallback:       lksdk.ParticipantCallback{},
		OnParticipantConnected:    callback.OnParticipantConnected,
		OnParticipantDisconnected: callback.OnParticipantDisconnected,
		OnDisconnected:            callback.OnDisconnected,
		OnReconnected:             callback.OnReconnected,
		OnReconnecting:            callback.OnReconnecting,
	}
	token, liveUrl, err = x.GetJoinToken(ctx, roomID, roomID)
	if err != nil {
		return "", "", "", errs.Wrap(err)
	}
	if _, err = lksdk.ConnectToRoomWithToken(x.conf.InnerURL, token, roomCallback); err != nil {
		return "", "", "", err
	}
	return room.Sid, token, liveUrl, nil
}

func (x *LiveKit) getLiveURL() string {
	if len(x.conf.URL) == 1 {
		return x.conf.URL[0]
	}
	return x.conf.URL[(atomic.AddUint64(&x.index, 1)-1)%uint64(len(x.conf.URL))]
}

func (x *LiveKit) RoomIsExist(ctx context.Context, roomID string) (string, error) {
	roomsResp, err := x.roomClient.ListRooms(ctx, &livekit.ListRoomsRequest{Names: []string{roomID}})
	if err != nil {
		return "", errs.Wrap(err)
	}
	if len(roomsResp.Rooms) > 0 {
		return roomsResp.Rooms[0].GetSid(), nil
	}
	return "", errs.ErrRecordNotFound.WrapMsg("roomIsNotExist")
}

func (x *LiveKit) GetRoomData(ctx context.Context, roomID string) (*meeting.MeetingMetadata, error) {
	resp, err := x.roomClient.ListRooms(ctx, &livekit.ListRoomsRequest{Names: []string{roomID}})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(resp.Rooms) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("roomIsNotExist")
	}
	var metaData meeting.MeetingMetadata
	if resp.Rooms[0].Metadata != "" {
		if err := json.Unmarshal([]byte(resp.Rooms[0].Metadata), &metaData); err != nil {
			return nil, errs.WrapMsg(err, "Unmarshal failed roomId:", roomID)
		}
		return &metaData, nil
	}
	return &metaData, nil
}

func (x *LiveKit) UpdateMetaData(ctx context.Context, updateData *meeting.MeetingMetadata) error {
	bytes, err := json.Marshal(&updateData)
	if err != nil {
		return errs.Wrap(err)
	}
	_, err = x.roomClient.UpdateRoomMetadata(ctx, &livekit.UpdateRoomMetadataRequest{
		Room:     updateData.Detail.Info.SystemGenerated.MeetingID,
		Metadata: string(bytes),
	})

	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func (x *LiveKit) CloseRoom(ctx context.Context, roomID string) error {
	_, err := x.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
		Room: roomID,
	})
	if err != nil {
		return errs.ErrInternalServer.WrapMsg(err.Error())
	}
	return nil
}

func (x *LiveKit) RemoveParticipant(ctx context.Context, roomID, userID string) error {
	_, err := x.roomClient.RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{Room: roomID, Identity: userID})
	if err != nil && !x.IsNotFound(err) {
		return err
	}
	return nil
}

func (x *LiveKit) IsNotFound(err error) bool {
	err = errs.Unwrap(err)
	if err == nil {
		return false
	}
	errCode, ok := err.(interface{ Code() twirp.ErrorCode })
	return ok && errCode.Code() == twirp.NotFound
}
