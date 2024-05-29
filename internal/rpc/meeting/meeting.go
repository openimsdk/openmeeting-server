package meeting

import (
	"context"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/constant"
	sysConstant "github.com/openimsdk/protocol/constant"

	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/idutil"
)

// BookMeeting Implement the MeetingServiceServer interface
func (s *meetingServer) BookMeeting(ctx context.Context, req *pbmeeting.BookMeetingReq) (*pbmeeting.BookMeetingResp, error) {
	resp := &pbmeeting.BookMeetingResp{}
	meetingDBInfo := &model.MeetingInfo{
		MeetingID:       idutil.OperationIDGenerator(),
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.Scheduled,
		CreatorUserID:   req.CreatorUserID,
	}

	err := s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}
	// fill in response data
	resp.Detail = s.generateRespSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	return resp, nil
}

func (s *meetingServer) CreateImmediateMeeting(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*pbmeeting.CreateImmediateMeetingResp, error) {
	resp := &pbmeeting.CreateImmediateMeetingResp{}
	meetingDBInfo := &model.MeetingInfo{
		MeetingID:       idutil.OperationIDGenerator(),
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.InProgress,
		CreatorUserID:   req.CreatorUserID,
	}
	_, token, liveUrl, err := s.meetingRtc.CreateRoom(ctx, meetingDBInfo.MeetingID)
	if err != nil {
		return resp, err
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}

	metaData := &pbmeeting.MeetingMetadata{}
	meetingDetail := s.generateRespSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	metaData.Detail = meetingDetail
	metaData.PersonalData = []*pbmeeting.PersonalData{s.getDefaultPersonalData(req.CreatorUserID)}
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

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get room data failed")
	}

	if req.Password != metaData.Detail.Info.CreatorDefinedMeeting.Password {
		return resp, errs.New("meeting password not match, please check and try again!")
	}

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.MeetingID)
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
		personalData := s.getDefaultPersonalData(req.UserID)
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

	// todo check user auth
	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.MeetingID)
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

	// Update the specific field based on the request
	updateData := map[string]any{
		"Title":           req.Title,
		"ScheduledTime":   req.ScheduledTime,
		"MeetingDuration": req.MeetingDuration,
		"Password":        req.Password,
	}
	metaData.Detail.Setting.CanParticipantsEnableCamera = req.CanParticipantsEnableCamera
	metaData.Detail.Setting.CanParticipantsUnmuteMicrophone = req.CanParticipantsUnmuteMicrophone
	metaData.Detail.Setting.CanParticipantsShareScreen = req.CanParticipantsShareScreen
	metaData.Detail.Setting.DisableCameraOnJoin = req.DisableCameraOnJoin
	metaData.Detail.Setting.DisableMicrophoneOnJoin = req.DisableMicrophoneOnJoin

	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, updateData); err != nil {
		return resp, err
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
	fmt.Println(metaData)
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
	if err := s.setParticipantPersonalSetting(ctx, metaData, req); err != nil {
		return resp, errs.WrapMsg(err, "set participant personal setting failed")
	}
	return resp, nil
}
