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

	//var req apistruct.SetPersonalSettingReq
	//if err := c.BindJSON(&req); err != nil {
	//	apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
	//	return
	//}
	//rpcReq := &meeting.SetPersonalMeetingSettingsReq{
	//	MeetingID: req.MeetingID,
	//	UserID:    req.UserID,
	//}
	//if req.Setting.CameraOnEntry != nil {
	//	rpcReq.CameraOnEntry = &pbwrapper.BoolValue{Value: *req.Setting.CameraOnEntry}
	//}
	//if req.Setting.MicrophoneOnEntry != nil {
	//	rpcReq.MicrophoneOnEntry = &pbwrapper.BoolValue{Value: *req.Setting.MicrophoneOnEntry}
	//}
	//resp, err := m.Client.SetPersonalMeetingSettings(c, rpcReq)
	//if err != nil {
	//	apiresp.GinError(c, err) // rpc call failed
	//	return
	//}
	//apiresp.GinSuccess(c, resp) // rpc call success}
}

func (m *MeetingApi) UpdateMeeting(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.UpdateMeeting, m.Client, c)

	//var req apistruct.UpdateMeetingReq

	//if err := c.BindJSON(&req); err != nil {
	//	apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
	//	return
	//}
	//
	//rpcReq := &meeting.UpdateMeetingRequest{
	//	MeetingID:      req.MeetingID,
	//	UpdatingUserID: req.UpdatingUserID,
	//}
	//if req.Title != nil {
	//	rpcReq.Title = &pbwrapper.StringValue{Value: *req.Title}
	//}
	//if req.Password != nil {
	//	rpcReq.Password = &pbwrapper.StringValue{Value: *req.Password}
	//}
	//if req.MeetingDuration != nil {
	//	rpcReq.MeetingDuration = &pbwrapper.Int64Value{Value: *req.MeetingDuration}
	//}
	//if req.ScheduledTime != nil {
	//	rpcReq.ScheduledTime = &pbwrapper.Int64Value{Value: *req.ScheduledTime}
	//}
	//
	//if req.CanParticipantsUnmuteMicrophone != nil {
	//	rpcReq.CanParticipantsUnmuteMicrophone = &pbwrapper.BoolValue{Value: *req.CanParticipantsUnmuteMicrophone}
	//}
	//if req.CanParticipantsEnableCamera != nil {
	//	rpcReq.CanParticipantsEnableCamera = &pbwrapper.BoolValue{Value: *req.CanParticipantsEnableCamera}
	//}
	//if req.DisableMicrophoneOnJoin != nil {
	//	rpcReq.DisableMicrophoneOnJoin = &pbwrapper.BoolValue{Value: *req.DisableMicrophoneOnJoin}
	//}
	//if req.CanParticipantsShareScreen != nil {
	//	rpcReq.CanParticipantsShareScreen = &pbwrapper.BoolValue{Value: *req.CanParticipantsShareScreen}
	//}
	//if req.DisableCameraOnJoin != nil {
	//	rpcReq.DisableCameraOnJoin = &pbwrapper.BoolValue{Value: *req.DisableCameraOnJoin}
	//}
	//
	//resp, err := m.Client.UpdateMeeting(c, rpcReq)
	//if err != nil {
	//	apiresp.GinError(c, err) // rpc call failed
	//	return
	//}
	//apiresp.GinSuccess(c, resp) // rpc call success
}

func (m *MeetingApi) OperateMeetingAllStream(c *gin.Context) {
	a2r.Call(meeting.MeetingServiceClient.OperateRoomAllStream, m.Client, c)
	//
	//var req apistruct.OperateMeetingAllStreamReq
	//
	//if err := c.BindJSON(&req); err != nil {
	//	apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
	//	return
	//}
	//rpcReq := &meeting.OperateRoomAllStreamReq{
	//	MeetingID:      req.MeetingID,
	//	OperatorUserID: req.OperatorUserID,
	//}
	//if req.CameraOnEntry != nil {
	//	rpcReq.CameraOnEntry = &pbwrapper.BoolValue{Value: *req.CameraOnEntry}
	//}
	//if req.MicrophoneOnEntry != nil {
	//	rpcReq.MicrophoneOnEntry = &pbwrapper.BoolValue{Value: *req.MicrophoneOnEntry}
	//}
	//resp, err := m.Client.OperateRoomAllStream(c, rpcReq)
	//if err != nil {
	//	apiresp.GinError(c, err) // rpc call failed
	//	return
	//}
	//apiresp.GinSuccess(c, resp) // rpc call success
}
