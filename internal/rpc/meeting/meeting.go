package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/constant"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	pbuser "github.com/openimsdk/openmeeting-server/pkg/protocol/user"
	sysConstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
)

// BookMeeting Implement the MeetingServiceServer interface
func (s *meetingServer) BookMeeting(ctx context.Context, req *pbmeeting.BookMeetingReq) (*pbmeeting.BookMeetingResp, error) {
	resp := &pbmeeting.BookMeetingResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.CreatorUserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	meetingDBInfo, err := s.generateMeetingDBData4Booking(ctx, req)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting data failed")
	}

	_, _, _, err = s.meetingRtc.CreateRoom(ctx, meetingDBInfo.MeetingID, req.CreatorUserID, nil)
	if err != nil {
		return resp, err
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}
	metaData := &pbmeeting.MeetingMetadata{}
	meetingDetail := s.generateClientRespMeetingSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	meetingDetail.Info.SystemGenerated.CreatorNickname = userInfo.Nickname
	metaData.Detail = meetingDetail
	metaData.PersonalData = []*pbmeeting.PersonalData{s.generateDefaultPersonalData(req.CreatorUserID)}
	// create meeting meta data
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	// fill in response data
	resp.Detail = s.generateClientRespMeetingSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	return resp, nil
}

func (s *meetingServer) CreateImmediateMeeting(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*pbmeeting.CreateImmediateMeetingResp, error) {
	resp := &pbmeeting.CreateImmediateMeetingResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.CreatorUserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	meetingDBInfo, err := s.generateMeetingDBData4Create(ctx, req)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting data failed")
	}

	participantMetaData := s.generateParticipantMetaData(userInfo)

	_, token, liveUrl, err := s.meetingRtc.CreateRoom(ctx, meetingDBInfo.MeetingID, req.CreatorUserID, participantMetaData)
	if err != nil {
		return resp, err
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}

	metaData := &pbmeeting.MeetingMetadata{}
	meetingDetail := s.generateClientRespMeetingSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	meetingDetail.Info.SystemGenerated.CreatorNickname = userInfo.Nickname
	metaData.Detail = meetingDetail
	metaData.PersonalData = []*pbmeeting.PersonalData{s.generateDefaultPersonalData(req.CreatorUserID)}
	// create meeting meta data
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	resp.Detail = meetingDetail
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) JoinMeeting(ctx context.Context, req *pbmeeting.JoinMeetingReq) (*pbmeeting.JoinMeetingResp, error) {
	resp := &pbmeeting.JoinMeetingResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.UserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get room data failed")
	}

	if req.Password != metaData.Detail.Info.CreatorDefinedMeeting.Password {
		return resp, errs.New("meeting password not match, please check and try again!")
	}

	participantMetaData := s.generateParticipantMetaData(userInfo)

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.UserID, participantMetaData)
	if err != nil {
		return resp, errs.WrapMsg(err, "get join token failed")
	}

	// update meta data to liveKit
	found := false
	for _, personalData := range metaData.PersonalData {
		if personalData.UserID == req.UserID {
			found = true
			break
		}
	}
	if !found {
		personalData := s.generateDefaultPersonalData(req.UserID)
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, errs.WrapMsg(err, "update meta data failed")
	}
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) GetMeetingToken(ctx context.Context, req *pbmeeting.GetMeetingTokenReq) (*pbmeeting.GetMeetingTokenResp, error) {
	resp := &pbmeeting.GetMeetingTokenResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.UserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	participantMetaData := s.generateParticipantMetaData(userInfo)

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.UserID, participantMetaData)
	if err != nil {
		return resp, err
	}

	resp.MeetingID = req.MeetingID
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) LeaveMeeting(ctx context.Context, req *pbmeeting.LeaveMeetingReq) (*pbmeeting.LeaveMeetingResp, error) {
	resp := &pbmeeting.LeaveMeetingResp{}

	if err := s.meetingRtc.RemoveParticipant(ctx, req.MeetingID, req.UserID); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *meetingServer) EndMeeting(ctx context.Context, req *pbmeeting.EndMeetingReq) (*pbmeeting.EndMeetingResp, error) {
	resp := &pbmeeting.EndMeetingResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}
	if !s.checkAuthPermission(metaData.Detail.Info.SystemGenerated.CreatorUserID, req.UserID) {
		return resp, errs.ErrArgs.WrapMsg("user did not have permission to end somebody's meeting")
	}
	// change status to completed
	updateData := map[string]any{
		"status": constant.Completed,
	}
	if err := s.meetingRtc.CloseRoom(ctx, req.MeetingID); err != nil {
		return resp, err
	}
	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, updateData); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *meetingServer) GetMeetings(ctx context.Context, req *pbmeeting.GetMeetingsReq) (*pbmeeting.GetMeetingsResp, error) {
	resp := &pbmeeting.GetMeetingsResp{}

	meetings, err := s.meetingStorageHandler.FindByStatus(ctx, req.Status)
	if err != nil {
		return resp, err
	}

	// Create response
	var meetingDetails []*pbmeeting.MeetingInfoSetting
	for _, meeting := range meetings {
		detailSetting, err := s.getMeetingDetailSetting(ctx, meeting)
		if err != nil {
			return resp, err
		}
		meetingDetails = append(meetingDetails, detailSetting)
	}
	resp.MeetingDetails = meetingDetails
	return resp, nil
}

func (s *meetingServer) GetMeeting(ctx context.Context, req *pbmeeting.GetMeetingReq) (*pbmeeting.GetMeetingResp, error) {
	resp := &pbmeeting.GetMeetingResp{}
	meetingDBInfo, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	detailSetting, err := s.getMeetingDetailSetting(ctx, meetingDBInfo)
	if err != nil {
		return resp, err
	}
	resp.MeetingDetail = detailSetting
	return resp, nil
}

func (s *meetingServer) UpdateMeeting(ctx context.Context, req *pbmeeting.UpdateMeetingRequest) (*pbmeeting.UpdateMeetingResp, error) {
	resp := &pbmeeting.UpdateMeetingResp{}

	_, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	updateData, liveKitUpdate := s.getUpdateData(metaData, req)

	if liveKitUpdate {
		if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
			return resp, err
		}
	}

	if len(*updateData) > 0 {
		if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, *updateData); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (s *meetingServer) GetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.GetPersonalMeetingSettingsReq) (*pbmeeting.GetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.GetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}
	for _, personalData := range metaData.PersonalData {
		if personalData.GetUserID() == req.UserID {
			resp.Setting = personalData.PersonalSetting
			break
		}
	}
	return resp, nil
}

func (s *meetingServer) SetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.SetPersonalMeetingSettingsReq) (*pbmeeting.SetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.SetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	result, err := common.GetKeyFromContext(ctx, sysConstant.OpUserID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get userid from context failed")
	}
	userID := result.(string)
	hostUser := s.getHostUserID(metaData)
	// non host user can not change other people's personal setting
	if userID != hostUser && userID != req.UserID {
		return resp, errs.New("do not have the permission to change other participant's setting")
	}

	if userID == req.UserID {
		// user set personal setting
		if err := s.setSelfPersonalSetting(ctx, metaData, req); err != nil {
			return resp, errs.WrapMsg(err, "set self personal setting failed")
		}
		return resp, nil
	}
	// admin set participant setting
	if err := s.setParticipantPersonalSetting(ctx, metaData, req); err != nil {
		return resp, errs.WrapMsg(err, "set participant personal setting failed")
	}
	return resp, nil
}

func (s *meetingServer) OperateRoomAllStream(ctx context.Context, req *pbmeeting.OperateRoomAllStreamReq) (*pbmeeting.OperateRoomAllStreamResp, error) {
	resp := &pbmeeting.OperateRoomAllStreamResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	hostUser := s.getHostUserID(metaData)
	if s.checkAuthPermission(hostUser, req.OperatorUserID) {
		return resp, errs.ErrNoPermission.WrapMsg("do not have the permission")
	}
	if req.MicrophoneOnEntry != nil {
		resp.StreamNotExistUserIDList, resp.FailedUserIDList, err = s.muteAllStream(ctx, req.MeetingID, audio, !req.MicrophoneOnEntry.Value)
		if err != nil {
			return resp, errs.WrapMsg(err, "operate room all microphone stream failed")
		}
	}

	if req.CameraOnEntry != nil {
		resp.StreamNotExistUserIDList, resp.StreamNotExistUserIDList, err = s.muteAllStream(ctx, req.MeetingID, video, !req.CameraOnEntry.Value)
		if err != nil {
			return resp, errs.WrapMsg(err, "operate room all camera stream failed")
		}
	}

	return resp, nil
}
