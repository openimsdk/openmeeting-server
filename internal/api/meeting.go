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
