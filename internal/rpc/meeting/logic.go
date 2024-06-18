package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/constant"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/pbwrapper"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/timeutil"
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
	if found {
		if req.CameraOnEntry == nil && req.MicrophoneOnEntry == nil {
			needUpdate = false
		}
	}
	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
	}
	if req.CameraOnEntry != nil {
		personalData.PersonalSetting.CameraOnEntry = req.CameraOnEntry.Value
	}
	if req.MicrophoneOnEntry != nil {
		personalData.PersonalSetting.MicrophoneOnEntry = req.MicrophoneOnEntry.Value
	}
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
	if found {
		if req.CameraOnEntry == nil && req.MicrophoneOnEntry == nil {
			needUpdate = false
		}
	}

	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
	}
	if req.CameraOnEntry != nil {
		personalData.LimitSetting.CameraOnEntry = req.CameraOnEntry.Value
	}
	if req.MicrophoneOnEntry != nil {
		personalData.LimitSetting.MicrophoneOnEntry = req.MicrophoneOnEntry.Value
	}
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

	if err := s.sendData(ctx, req.MeetingID, req.UserID, req.CameraOnEntry, req.MicrophoneOnEntry); err != nil {
		return errs.WrapMsg(err, "send data failed")
	}

	return nil
}

func (s *meetingServer) sendData(ctx context.Context, roomID, userID string, cameraOn, microphoneOn *pbwrapper.BoolValue) error {
	operationData := &pbmeeting.UserOperationData{
		UserID: userID,
	}
	if cameraOn != nil {
		operationData.CameraOnEntry = cameraOn.Value
	}
	if microphoneOn != nil {
		operationData.MicrophoneOnEntry = microphoneOn.Value
	}

	sendData := &pbmeeting.StreamOperateData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		Operation:      []*pbmeeting.UserOperationData{operationData},
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

func (s *meetingServer) send2AllParticipant(ctx context.Context, req *pbmeeting.OperateRoomAllStreamReq, StreamNotExistUserIDList, FailedUserIDList []string) error {
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
	sendData := &pbmeeting.StreamOperateData{
		OperatorUserID: mcontext.GetOpUserID(ctx),
		Operation:      operationList,
	}

	if err := s.meetingRtc.SendRoomData(ctx, req.MeetingID, nil, sendData); err != nil {
		return errs.WrapMsg(err, "send room data failed")
	}

	return nil
}

func (s *meetingServer) refreshMeetingStatus(ctx context.Context) {
	meetings, err := s.meetingStorageHandler.FindByStatus(ctx, []string{constant.InProgress, constant.Scheduled})
	if err != nil {
		log.ZError(ctx, "find meetings failed", err)
		return
	}
	nowTimestamp := timeutil.GetCurrentTimestampBySecond()
	for _, one := range meetings {
		if one.StartTime+one.MeetingDuration < nowTimestamp {
			updateData := map[string]any{
				"status": constant.Completed,
			}
			if err := s.meetingStorageHandler.Update(ctx, one.MeetingID, updateData); err != nil {
				log.ZError(ctx, "update meeting status failed", err)
			}
		} else if one.StartTime+one.MeetingDuration < nowTimestamp && one.StartTime > nowTimestamp {
			if one.Status == constant.InProgress {
				continue
			}
			updateData := map[string]any{
				"status": constant.InProgress,
			}
			if err := s.meetingStorageHandler.Update(ctx, one.MeetingID, updateData); err != nil {
				log.ZError(ctx, "update meeting status failed", err)
			}
		}
	}
}
