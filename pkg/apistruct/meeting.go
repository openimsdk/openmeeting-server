package apistruct

type UpdateMeetingReq struct {
	MeetingID                       string  `json:"meetingID"`
	UpdatingUserID                  string  `json:"updatingUserID"`
	Title                           *string `json:"title"`
	ScheduledTime                   *int64  `json:"scheduledTime"`
	MeetingDuration                 *int64  `json:"meetingDuration"`
	Password                        *string `json:"password"`
	CanParticipantsEnableCamera     *bool   `json:"canParticipantsEnableCamera"`
	CanParticipantsUnmuteMicrophone *bool   `json:"canParticipantsUnmuteMicrophone"`
	CanParticipantsShareScreen      *bool   `json:"canParticipantsShareScreen"`
	DisableCameraOnJoin             *bool   `json:"disableCameraOnJoin"`
	DisableMicrophoneOnJoin         *bool   `json:"disableMicrophoneOnJoin"`
}

type OperateMeetingAllStreamReq struct {
	MeetingID         string `json:"meetingID"`
	OperatorUserID    string `json:"operatorUserID"`
	MicrophoneOnEntry *bool  `json:"canParticipantsEnableCamera"`
	CameraOnEntry     *bool  `json:"canParticipantsUnmuteMicrophone"`
}