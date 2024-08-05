package livekit

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/rtc"
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
	log.ZDebug(ctx, "OnRoomParticipantConnected", "roomID:", r.roomID, "userID:", userID)
	// set default host when the first one coming in
	metaData, err := r.liveKit.GetRoomData(ctx, r.roomID)
	if err != nil {
		return
	}
	hostUserID := metaData.Detail.Info.CreatorDefinedMeeting.HostUserID

	participants, err := r.liveKit.ListParticipants(ctx, r.roomID)
	if err != nil {
		return
	}
	log.ZDebug(ctx, "OnRoomParticipantConnected",
		"room participant number:", len(participants),
		"hostID", hostUserID)

	// when first coming delete auto change host
	if len(participants) == 2 {
		// when first comer is creator, he is not host, so set him as the host
		if userID == metaData.Detail.Info.SystemGenerated.CreatorUserID && userID != hostUserID {
			metaData.Detail.Info.CreatorDefinedMeeting.HostUserID = userID
			log.CInfo(ctx, "set host info as default when creator is the first one to come in",
				"roomID:", r.roomID, "new host:", metaData.Detail.Info.SystemGenerated.CreatorUserID)
			if err := r.liveKit.UpdateMetaData(ctx, metaData); err != nil {
				log.ZError(ctx, "update meta room data change host info failed", err,
					"new host:", metaData.Detail.Info.SystemGenerated.CreatorUserID)
			}
		}
	}
	if hostUserID == "" {
		metaData.Detail.Info.CreatorDefinedMeeting.HostUserID = metaData.Detail.Info.SystemGenerated.CreatorUserID
		log.CInfo(ctx, "set host info as default when last host is nil",
			"roomID:", r.roomID, "new host:", metaData.Detail.Info.SystemGenerated.CreatorUserID)
		if err := r.liveKit.UpdateMetaData(ctx, metaData); err != nil {
			log.ZError(ctx, "update meta room data change host info failed", err,
				"new host:", metaData.Detail.Info.SystemGenerated.CreatorUserID)
		}
	}
}

func (r *CallbackLiveKit) OnRoomParticipantDisconnected(ctx context.Context, userID string) {
	log.ZWarn(ctx, "OnRoomParticipantDisconnected", nil, "userID:", userID)
	if err := r.liveKit.RemoveParticipant(ctx, r.roomID, userID); err != nil {
		log.ZWarn(ctx, "remove participant failed", err)
	}
	// auto change host to creator
	metaData, err := r.liveKit.GetRoomData(ctx, r.roomID)
	if err != nil {
		return
	}
	hostUserID := metaData.Detail.Info.CreatorDefinedMeeting.HostUserID
	creatorUserID := metaData.Detail.Info.SystemGenerated.CreatorUserID
	if hostUserID == userID && creatorUserID != hostUserID {
		log.CInfo(ctx, "change host info when last host disconnected", "roomID:", r.roomID, "old host:", hostUserID, "default host:", creatorUserID)
		metaData.Detail.Info.CreatorDefinedMeeting.HostUserID = creatorUserID
		if err := r.liveKit.UpdateMetaData(ctx, metaData); err != nil {
			log.ZError(ctx, "update meta room data change host info failed", err, "old host:", hostUserID, "default host:", creatorUserID)
		}
	}
}

func (r *CallbackLiveKit) OnRoomDisconnected(ctx context.Context) {
	log.ZWarn(ctx, "OnRoomDisconnected", nil, "roomID", r.roomID)
	participants, err := r.liveKit.ListParticipants(ctx, r.roomID)
	if err != nil {
		log.ZWarn(ctx, "remove participant failed", err, r.roomID)
		return
	}
	for _, p := range participants {
		if err := r.liveKit.RemoveParticipant(ctx, r.roomID, p.Identity); err != nil {
			log.ZWarn(ctx, "remove participant failed", err, p.Identity)
		}
	}
}

func (r *CallbackLiveKit) OnMeetingDisconnected(ctx context.Context, roomID string) {
	log.ZWarn(ctx, "OnMeetingDisconnected", nil, roomID)
}

func (r *CallbackLiveKit) OnMeetingUnmute(ctx context.Context, roomID string, streamType string, mute bool, userIDs []string) {
	log.ZWarn(ctx, "OnMeetingUnmute", nil, roomID)
}
