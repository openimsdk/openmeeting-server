package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/errs"
)

const (
	video = "video"
	audio = "audio"
)

func (s *meetingServer) getHostUserID(metadata *pbmeeting.MeetingMetadata) string {
	return metadata.Detail.Info.SystemGenerated.CreatorUserID
}

func (s *meetingServer) checkAuthPermission(hostUserID, requestUserID string) bool {
	return hostUserID == requestUserID
}

func (s *meetingServer) getMeetingDetailSetting(ctx context.Context, info *model.MeetingInfo) (*pbmeeting.MeetingInfoSetting, error) {
	// Fill in response data
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID: info.CreatorUserID,
		Status:        info.Status,
		StartTime:     info.StartTime,
		MeetingID:     info.MeetingID,
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
	}

	return meetingInfoSetting, nil
}

// generateMeetingInfoSetting generates MeetingInfoSetting from the given request and meeting ID.
func (s *meetingServer) generateRespSetting(
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

func (s *meetingServer) getDefaultPersonalData(userID string) *pbmeeting.PersonalData {
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

func (s *meetingServer) setSelfPersonalSetting(ctx context.Context, metaData *pbmeeting.MeetingMetadata, req *pbmeeting.SetPersonalMeetingSettingsReq) error {
	found := false
	var personalData *pbmeeting.PersonalData
	for _, one := range metaData.PersonalData {
		if one.GetUserID() == req.UserID {
			personalData = one
			found = true
			break
		}
	}
	needUpdate := true
	if found && personalData.PersonalSetting.CameraOnEntry == req.Setting.CameraOnEntry &&
		personalData.PersonalSetting.MicrophoneOnEntry == req.Setting.MicrophoneOnEntry {
		needUpdate = false
	}

	personalData.PersonalSetting = req.Setting
	toggle := s.checkUserEnableCamera(metaData.Detail.Setting, personalData)
	if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, video, !toggle); err != nil {
		return errs.WrapMsg(err, "toggle camera stream failed")
	}

	toggle = s.checkUserEnableMicrophone(metaData.Detail.Setting, personalData)
	if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, audio, !toggle); err != nil {
		return errs.WrapMsg(err, "toggle microphone stream failed")
	}

	// judge whether user need to change or not
	if !found {
		metaData.PersonalData = append(metaData.PersonalData, &pbmeeting.PersonalData{
			UserID:          req.UserID,
			PersonalSetting: req.Setting,
			LimitSetting: &pbmeeting.PersonalMeetingSetting{
				CameraOnEntry:     true,
				MicrophoneOnEntry: true,
			},
		})
	}
	if !needUpdate {
		// no need update
		return nil
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return errs.WrapMsg(err, "update meta data failed")
	}

	return nil
}

// setParticipantPersonalSetting set setting
func (s *meetingServer) setParticipantPersonalSetting(ctx context.Context, metaData *pbmeeting.MeetingMetadata, req *pbmeeting.SetPersonalMeetingSettingsReq) error {
	found := false
	var personalData *pbmeeting.PersonalData
	for _, one := range metaData.PersonalData {
		if one.GetUserID() == req.UserID {
			personalData = one
			found = true
			break
		}
	}

	// judge whether user need to change or not
	needUpdate := true
	if found && personalData.LimitSetting.MicrophoneOnEntry == req.Setting.MicrophoneOnEntry &&
		personalData.LimitSetting.CameraOnEntry == req.Setting.CameraOnEntry {
		needUpdate = false
	}

	personalData.LimitSetting = req.Setting
	if !s.checkUserEnableCamera(metaData.Detail.Setting, personalData) {
		// no need to care the scene that turning on the camera
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, video, true); err != nil {
			return errs.WrapMsg(err, "toggle camera stream failed")
		}
	}

	if !s.checkUserEnableMicrophone(metaData.Detail.Setting, personalData) {
		// no need to care the scene that turning on the microphone
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, audio, true); err != nil {
			return errs.WrapMsg(err, "toggle microphone stream failed")
		}
	}

	if !found {
		metaData.PersonalData = append(metaData.PersonalData, &pbmeeting.PersonalData{
			UserID:          req.UserID,
			PersonalSetting: req.Setting,
		})
	}
	if !needUpdate {
		// no need to update
		return nil
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return errs.WrapMsg(err, "update meta data failed")
	}
	return nil
}
