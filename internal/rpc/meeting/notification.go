package meeting

import (
	"context"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	pbwrapper "github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func (s *meetingServer) sendMeetingHostData2Client(ctx context.Context, roomID, operateUserID, userID, hostType string) error {
	userInfo, err := s.userRpc.GetUserInfo(ctx, operateUserID)
	if err != nil {
		return errs.WrapMsg(err, "get user info failed")
	}

	hostInfo := &pbmeeting.MeetingHostData{
		OperatorNickname: userInfo.Nickname,
		UserID:           userID,
		HostType:         hostType,
	}

	sendData := &pbmeeting.NotifyMeetingData{
		OperatorUserID: operateUserID,
		MessageType:    &pbmeeting.NotifyMeetingData_MeetingHostData{MeetingHostData: hostInfo},
	}

	if err := s.meetingRtc.SendRoomData(ctx, roomID, &[]string{userID}, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}
	return nil
}

func (s *meetingServer) sendStreamOperateData2Client(ctx context.Context, roomID, userID string, cameraOn, microphoneOn *pbwrapper.BoolValue) error {
	operationData := &pbmeeting.UserOperationData{
		UserID: userID,
	}
	operationData.CameraOnEntry = cameraOn
	operationData.MicrophoneOnEntry = microphoneOn

	streamOperationData := &pbmeeting.StreamOperateData{
		Operation: []*pbmeeting.UserOperationData{operationData},
	}

	sendData := &pbmeeting.NotifyMeetingData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		MessageType:    &pbmeeting.NotifyMeetingData_StreamOperateData{StreamOperateData: streamOperationData},
	}

	if err := s.meetingRtc.SendRoomData(ctx, roomID, &[]string{userID}, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}
	return nil
}

func (s *meetingServer) broadcastStreamOperateData(ctx context.Context, req *pbmeeting.OperateRoomAllStreamReq, StreamNotExistUserIDList, FailedUserIDList []string) error {
	userIDs, err := s.meetingRtc.GetParticipantUserIDs(ctx, req.MeetingID)
	if err != nil {
		return err
	}
	userMap := make(map[string]bool)
	var setUserIDs []string
	for _, v := range userIDs {
		userMap[v] = true
	}
	for _, v := range StreamNotExistUserIDList {
		if _, ok := userMap[v]; !ok {
			delete(userMap, v)
		}
	}
	for _, v := range FailedUserIDList {
		if _, ok := userMap[v]; !ok {
			delete(userMap, v)
		}
	}
	for k, _ := range userMap {
		setUserIDs = append(setUserIDs, k)
	}

	var operationList []*pbmeeting.UserOperationData
	for _, v := range setUserIDs {
		operationData := &pbmeeting.UserOperationData{
			UserID: v,
		}
		operationData.CameraOnEntry = req.CameraOnEntry
		operationData.MicrophoneOnEntry = req.MicrophoneOnEntry
		operationList = append(operationList, operationData)
	}

	streamOperationData := &pbmeeting.StreamOperateData{
		Operation: operationList,
	}

	sendData := &pbmeeting.NotifyMeetingData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		MessageType:    &pbmeeting.NotifyMeetingData_StreamOperateData{StreamOperateData: streamOperationData},
	}

	if err := s.meetingRtc.SendRoomData(ctx, req.MeetingID, nil, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}

	return nil
}

func (s *meetingServer) notifyKickOffMeetingInfo2Client(ctx context.Context, roomID, UserID, Reason string, reasonCode pbmeeting.KickOffReason) error {
	userInfo, err := s.userRpc.GetUserInfo(ctx, UserID)
	if err != nil {
		return errs.WrapMsg(err, "get user info failed")
	}

	kickOffData := &pbmeeting.KickOffMeetingData{
		UserID:     UserID,
		Nickname:   userInfo.Nickname,
		ReasonCode: reasonCode,
		Reason:     Reason,
	}

	sendData := &pbmeeting.NotifyMeetingData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		MessageType:    &pbmeeting.NotifyMeetingData_KickOffMeetingData{KickOffMeetingData: kickOffData},
	}
	if err := s.meetingRtc.SendRoomData(ctx, roomID, nil, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}

	return nil
}
