package meeting

import (
	"context"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	"github.com/openimsdk/tools/errs"
)

func (s *meetingServer) checkAuthPermission(creatorUserID, hostUserID, requestUserID string) bool {
	return hostUserID == requestUserID || creatorUserID == requestUserID
}

func (s *meetingServer) checkUserEnableCamera(setting *pbmeeting.MeetingSetting, personalData *pbmeeting.PersonalData) bool {
	if setting.CanParticipantsEnableCamera && personalData.PersonalSetting.CameraOnEntry && personalData.LimitSetting.CameraOnEntry {
		return true
	}
	return false
}

func (s *meetingServer) checkUserEnableMicrophone(setting *pbmeeting.MeetingSetting, personalData *pbmeeting.PersonalData) bool {
	// only when three condition enable can turn on the microphone
	if setting.CanParticipantsUnmuteMicrophone && personalData.PersonalSetting.MicrophoneOnEntry && personalData.LimitSetting.MicrophoneOnEntry {
		return true
	}
	return false
}

func (s *meetingServer) checkUserInMeeting(ctx context.Context, userID string) (bool, error) {
	rooms, err := s.meetingRtc.GetAllRooms(ctx)
	if err != nil {
		return true, err
	}

	for _, room := range rooms {
		userIDs, err := s.meetingRtc.GetParticipantUserIDs(ctx, room.Name)
		if err != nil {
			return true, errs.WrapMsg(err, "get participants failed")
		}
		//check if user is already in meeting
		for _, userIdentity := range userIDs {
			if userIdentity == userID {
				return true, nil
			}
		}
	}
	return false, nil
}
