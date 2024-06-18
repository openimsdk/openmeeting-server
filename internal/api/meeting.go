package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/tools/a2r"
)

type MeetingApi rpcclient.Meeting

func NewMeetingApi(client rpcclient.Meeting) MeetingApi {
	return MeetingApi(client)
}

func (m *MeetingApi) BookMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.BookMeeting, m.Client, c)
}

func (m *MeetingApi) CreateImmediateMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.CreateImmediateMeeting, m.Client, c)
}

func (m *MeetingApi) JoinMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.JoinMeeting, m.Client, c)
}

func (m *MeetingApi) GetMeetingToken(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.GetMeetingToken, m.Client, c)
}

func (m *MeetingApi) LeaveMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.LeaveMeeting, m.Client, c)
}

func (m *MeetingApi) EndMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.EndMeeting, m.Client, c)
}

func (m *MeetingApi) GetMeetings(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.GetMeetings, m.Client, c)
}

func (m *MeetingApi) GetMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.GetMeeting, m.Client, c)
}

func (m *MeetingApi) GetPersonalMeetingSettings(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.GetPersonalMeetingSettings, m.Client, c)
}

func (m *MeetingApi) ModifyMeetingParticipantNickName(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.ModifyMeetingParticipantNickName, m.Client, c)
}

func (m *MeetingApi) RemoveMeetingParticipants(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.RemoveParticipants, m.Client, c)
}

func (m *MeetingApi) SetMeetingHostInfo(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.SetMeetingHostInfo, m.Client, c)
}

func (m *MeetingApi) SetPersonalMeetingSettings(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.SetPersonalMeetingSettings, m.Client, c)
}

func (m *MeetingApi) UpdateMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.UpdateMeeting, m.Client, c)
}

func (m *MeetingApi) OperateMeetingAllStream(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.OperateRoomAllStream, m.Client, c)
}
