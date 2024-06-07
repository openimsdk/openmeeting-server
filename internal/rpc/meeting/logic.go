package meeting

import (
	"context"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

const (
	video = "video"
	audio = "audio"
)

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
	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
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
		metaData.PersonalData = append(metaData.PersonalData, personalData)
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

	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
	}
	personalData.LimitSetting = req.Setting

	var (
		cameraOn     = false
		microphoneOn = false
	)
	if cameraOn = s.checkUserEnableCamera(metaData.Detail.Setting, personalData); !cameraOn {
		// no need to care the scene that turning on the camera
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, video, true); err != nil {
			return errs.WrapMsg(err, "toggle camera stream failed")
		}
	}

	if microphoneOn = s.checkUserEnableMicrophone(metaData.Detail.Setting, personalData); !microphoneOn {
		// no need to care the scene that turning on the microphone
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, audio, true); err != nil {
			return errs.WrapMsg(err, "toggle microphone stream failed")
		}
	}

	if !found {
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}
	if !needUpdate {
		// no need to update
		return nil
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return errs.WrapMsg(err, "update meta data failed")
	}

	if err := s.sendData(ctx, req.MeetingID, req.UserID, req.Setting.CameraOnEntry, req.Setting.MicrophoneOnEntry); err != nil {
		return errs.WrapMsg(err, "send data failed")
	}

	return nil
}

func (s *meetingServer) sendData(ctx context.Context, roomID, userID string, cameraOn, microphoneOn bool) error {
	sendData := &pbmeeting.StreamOperateData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		Operation: []*pbmeeting.UserOperationData{&pbmeeting.UserOperationData{
			UserID: userID,
			Operation: &pbmeeting.PersonalMeetingSetting{
				CameraOnEntry:     cameraOn,
				MicrophoneOnEntry: microphoneOn,
			},
		}},
	}

	if err := s.meetingRtc.SendRoomData(ctx, roomID, &[]string{userID}, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}
	return nil
}

func (s *meetingServer) muteAllStream(ctx context.Context, roomID, streamType string, mute bool) (streamNotExistUserIDList []string, failedUserIDList []string, err error) {
	participants, err := s.meetingRtc.ListParticipants(ctx, roomID)
	if err != nil {
		return nil, nil, errs.WrapMsg(err, "get participant list failed")
	}
	for _, v := range participants {
		err := s.meetingRtc.ToggleMimeStream(ctx, roomID, v.Identity, streamType, mute)
		if err != nil {
			log.ZError(ctx, "muteAllStream failed", err)
			if errs.ErrRecordNotFound.Is(err) {
				streamNotExistUserIDList = append(streamNotExistUserIDList, v.Identity)
			} else {
				failedUserIDList = append(failedUserIDList, v.Identity)
			}
		}
	}
	return
}
