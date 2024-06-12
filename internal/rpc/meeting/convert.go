package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/constant"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	pbuser "github.com/openimsdk/openmeeting-server/pkg/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/timeutil"
)

func (s *meetingServer) getHostUserID(metadata *pbmeeting.MeetingMetadata) string {
	return metadata.Detail.Info.SystemGenerated.CreatorUserID
}

func (s *meetingServer) generateMeetingDBData4Booking(ctx context.Context, req *pbmeeting.BookMeetingReq) (*model.MeetingInfo, error) {
	meetingID, err := s.meetingStorageHandler.GenerateMeetingID(ctx)
	if err != nil {
		return nil, errs.WrapMsg(err, "generate meeting id failed")
	}

	return &model.MeetingInfo{
		MeetingID:       meetingID,
		Title:           req.CreatorDefinedMeetingInfo.Title,
		StartTime:       req.CreatorDefinedMeetingInfo.ScheduledTime,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.Scheduled,
		CreatorUserID:   req.CreatorUserID,
	}, nil
}

func (s *meetingServer) generateMeetingDBData4Create(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*model.MeetingInfo, error) {
	meetingID, err := s.meetingStorageHandler.GenerateMeetingID(ctx)
	if err != nil {
		return nil, errs.WrapMsg(err, "generate meeting id failed")
	}

	return &model.MeetingInfo{
		MeetingID:       meetingID,
		Title:           req.CreatorDefinedMeetingInfo.Title,
		StartTime:       timeutil.GetCurrentTimestampBySecond(),
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.InProgress,
		CreatorUserID:   req.CreatorUserID,
	}, nil
}

func (s *meetingServer) generateParticipantMetaData(userInfo *pbuser.GetUserInfoResp) *pbmeeting.ParticipantMetaData {
	return &pbmeeting.ParticipantMetaData{
		UserInfo: &pbmeeting.UserInfo{
			UserID:   userInfo.UserID,
			Nickname: userInfo.Nickname,
			Account:  userInfo.Account,
		},
	}
}

func (s *meetingServer) generateDefaultPersonalData(userID string) *pbmeeting.PersonalData {
	return &pbmeeting.PersonalData{
		UserID: userID,
		PersonalSetting: &pbmeeting.PersonalMeetingSetting{
			CameraOnEntry:     false,
			MicrophoneOnEntry: false,
		},
		LimitSetting: &pbmeeting.PersonalMeetingSetting{
			CameraOnEntry:     true,
			MicrophoneOnEntry: true,
		},
	}
}

func (s *meetingServer) getMeetingDetailSetting(ctx context.Context, info *model.MeetingInfo) (*pbmeeting.MeetingInfoSetting, error) {
	// Fill in response data
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: info.CreatorUserID})
	if err != nil {
		return nil, errs.WrapMsg(err, "get user info failed")
	}

	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID:   info.CreatorUserID,
		Status:          info.Status,
		StartTime:       info.StartTime,
		MeetingID:       info.MeetingID,
		CreatorNickname: userInfo.Nickname,
	}
	creatorInfo := &pbmeeting.CreatorDefinedMeetingInfo{
		Title:           info.Title,
		ScheduledTime:   info.ScheduledTime,
		MeetingDuration: info.MeetingDuration,
		Password:        info.Password,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info: meetingInfo,
	}
	metaData, err := s.meetingRtc.GetRoomData(ctx, info.MeetingID)
	if err == nil {
		meetingInfoSetting.Setting = metaData.Detail.Setting
		meetingInfoSetting.Info.SystemGenerated.CreatorNickname = metaData.Detail.Info.SystemGenerated.CreatorNickname
	}

	return meetingInfoSetting, nil
}

// generateMeetingInfoSetting generates MeetingInfoSetting from the given request and meeting ID.
func (s *meetingServer) generateClientRespMeetingSetting(
	meetingSetting *pbmeeting.MeetingSetting,
	defineMeetingInfo *pbmeeting.CreatorDefinedMeetingInfo, meeting *model.MeetingInfo) *pbmeeting.MeetingInfoSetting {
	// Fill in response data
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID: meeting.CreatorUserID,
		Status:        meeting.Status,
		StartTime:     meeting.ScheduledTime, // Scheduled start time as the actual start time
		MeetingID:     meeting.MeetingID,
	}
	// Combine system-generated and creator-defined info
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: defineMeetingInfo,
	}
	// Create MeetingInfoSetting
	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info:    meetingInfo,
		Setting: meetingSetting,
	}
	return meetingInfoSetting
}

func (s *meetingServer) getUpdateData(metaData *pbmeeting.MeetingMetadata, req *pbmeeting.UpdateMeetingRequest) (*map[string]any, bool) {
	// Update the specific field based on the request
	liveKitUpdate := false
	updateData := map[string]any{}
	if req.Title != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.Title = req.Title.Value
		updateData["Title"] = req.Title.Value
	}
	if req.ScheduledTime != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.ScheduledTime = req.ScheduledTime.Value
		updateData["ScheduledTime"] = req.ScheduledTime.Value
	}
	if req.MeetingDuration != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.MeetingDuration = req.ScheduledTime.Value
		updateData["MeetingDuration"] = req.MeetingDuration.Value
	}
	if req.Password != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.Password = req.Password.Value
		updateData["Password"] = req.Password.Value
	}

	if req.CanParticipantsEnableCamera != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.CanParticipantsEnableCamera = req.CanParticipantsEnableCamera.Value
	}
	if req.CanParticipantsUnmuteMicrophone != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.CanParticipantsUnmuteMicrophone = req.CanParticipantsUnmuteMicrophone.Value
	}
	if req.CanParticipantsShareScreen != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.CanParticipantsShareScreen = req.CanParticipantsShareScreen.Value
	}
	if req.DisableCameraOnJoin != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.DisableCameraOnJoin = req.DisableCameraOnJoin.Value
	}
	if req.DisableMicrophoneOnJoin != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.DisableMicrophoneOnJoin = req.DisableMicrophoneOnJoin.Value
	}
	return &updateData, liveKitUpdate
}
