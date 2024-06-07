package meeting

import pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"

func (s *meetingServer) checkAuthPermission(hostUserID, requestUserID string) bool {
	return hostUserID == requestUserID
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
