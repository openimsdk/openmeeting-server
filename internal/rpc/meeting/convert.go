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

	dbInfo := &model.MeetingInfo{
		MeetingID:       meetingID,
		Title:           req.CreatorDefinedMeetingInfo.Title,
		StartTime:       req.CreatorDefinedMeetingInfo.ScheduledTime,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		TimeZone:        req.CreatorDefinedMeetingInfo.TimeZone,
		Status:          constant.Scheduled,
		EndDate:         req.RepeatInfo.EndDate,
		RepeatType:      req.RepeatInfo.RepeatType,
		CreatorUserID:   req.CreatorUserID,
	}
	if req.RepeatInfo.RepeatType == constant.RepeatCustom {
		dbInfo.UintType = req.RepeatInfo.UintType
		dbInfo.Interval = req.RepeatInfo.Interval
	}

	return dbInfo, nil
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

func (s *meetingServer) generateMeetingMetaData(ctx context.Context, req *pbmeeting.BookMeetingReq, info *model.MeetingInfo) (*pbmeeting.MeetingMetadata, error) {
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: info.CreatorUserID})
	if err != nil {
		return nil, errs.WrapMsg(err, "get user info failed")
	}

	metaData := &pbmeeting.MeetingMetadata{}
	metaData.PersonalData = []*pbmeeting.PersonalData{s.generateDefaultPersonalData(req.CreatorUserID)}
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID:   info.CreatorUserID,
		Status:          info.Status,
		StartTime:       info.StartTime,
		MeetingID:       info.MeetingID,
		CreatorNickname: userInfo.Nickname,
	}
	creatorInfo := &pbmeeting.CreatorDefinedMeetingInfo{
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		TimeZone:        req.CreatorDefinedMeetingInfo.TimeZone,
		HostUserID:      req.CreatorUserID,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	metaData.Detail = &pbmeeting.MeetingInfoSetting{
		Setting:    req.Setting,
		Info:       meetingInfo,
		RepeatInfo: req.RepeatInfo,
	}
	return metaData, nil
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
		TimeZone:        info.TimeZone,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	repeatInfo := &pbmeeting.MeetingRepeatInfo{
		EndDate:    info.EndDate,
		RepeatType: info.RepeatType,
		UintType:   info.UintType,
		Interval:   info.Interval,
	}

	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info:       meetingInfo,
		RepeatInfo: repeatInfo,
	}
	metaData, err := s.meetingRtc.GetRoomData(ctx, info.MeetingID)
	if err == nil {
		meetingInfoSetting.Setting = metaData.Detail.Setting
		meetingInfoSetting.Info.SystemGenerated.CreatorNickname = metaData.Detail.Info.SystemGenerated.CreatorNickname
		meetingInfoSetting.Info.CreatorDefinedMeeting.MeetingDuration = metaData.Detail.Info.CreatorDefinedMeeting.MeetingDuration
		meetingInfoSetting.Info.CreatorDefinedMeeting.HostUserID = metaData.Detail.Info.CreatorDefinedMeeting.HostUserID
		meetingInfoSetting.Info.CreatorDefinedMeeting.CoHostUSerID = metaData.Detail.Info.CreatorDefinedMeeting.CoHostUSerID
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
		updateData["title"] = req.Title.Value
	}
	if req.ScheduledTime != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.ScheduledTime = req.ScheduledTime.Value
		updateData["scheduled_time"] = req.ScheduledTime.Value
	}
	if req.MeetingDuration != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.MeetingDuration = req.MeetingDuration.Value
		updateData["meeting_duration"] = req.MeetingDuration.Value
	}
	if req.Password != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.Password = req.Password.Value
		updateData["password"] = req.Password.Value
	}

	if req.RepeatInfo != nil {
		metaData.Detail.RepeatInfo = req.RepeatInfo
		updateData["repeat_type"] = req.RepeatInfo.RepeatType
		if req.RepeatInfo.RepeatType == constant.RepeatCustom {
			updateData["uint_type"] = req.RepeatInfo.UintType
			updateData["interval"] = req.RepeatInfo.Interval
		}
	}

	if req.TimeZone != nil {
		liveKitUpdate = true
		metaData.Detail.Info.CreatorDefinedMeeting.TimeZone = req.TimeZone.Value
		updateData["time_zone"] = req.TimeZone.Value
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
	if req.CanParticipantJoinMeetingEarly != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.CanParticipantJoinMeetingEarly = req.CanParticipantJoinMeetingEarly.Value
	}
	if req.AudioEncouragement != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.AudioEncouragement = req.AudioEncouragement.Value
	}
	if req.LockMeeting != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.LockMeeting = req.LockMeeting.Value
	}
	if req.VideoMirroring != nil {
		liveKitUpdate = true
		metaData.Detail.Setting.VideoMirroring = req.VideoMirroring.Value
	}

	return &updateData, liveKitUpdate
}

func (s *meetingServer) mergeAndUnique(array1, array2 []string) []string {
	exists := make(map[string]bool)
	var result []string

	for _, v := range array1 {
		if !exists[v] {
			exists[v] = true
			result = append(result, v)
		}
	}
	for _, v := range array2 {
		if !exists[v] {
			exists[v] = true
			result = append(result, v)
		}
	}
	return result
}
