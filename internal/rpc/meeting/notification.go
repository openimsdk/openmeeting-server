package meeting

import (
	"context"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	pbuser "github.com/openimsdk/protocol/openmeeting/user"
	pbwrapper "github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func (s *meetingServer) sendMeetingHostData2Client(ctx context.Context, roomID, operateUserID, userID, hostType string) error {
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: operateUserID})
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
	if cameraOn != nil {
		operationData.CameraOnEntry = cameraOn.Value
	}
	if microphoneOn != nil {
		operationData.MicrophoneOnEntry = microphoneOn.Value
	}
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
		if req.CameraOnEntry != nil {
			operationData.CameraOnEntry = req.CameraOnEntry.Value
		}
		if req.MicrophoneOnEntry != nil {
			operationData.MicrophoneOnEntry = req.MicrophoneOnEntry.Value
		}
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
