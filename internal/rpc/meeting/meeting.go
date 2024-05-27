package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting/rtc"
	"github.com/openimsdk/openmeeting-server/internal/rpc/meeting/rtc/livekit"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	mgo2 "github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/constant"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	registry "github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/idutil"
	"google.golang.org/grpc"
)

type meetingServer struct {
	meetingStorageHandler controller.Meeting
	RegisterCenter        registry.SvcDiscoveryRegistry
	meetingRtc            rtc.MeetingRtc
	config                *Config
}

type Config struct {
	Rpc       config.Meeting
	Redis     config.Redis
	Mongo     config.Mongo
	Discovery config.Discovery
	Share     config.Share
	Rtc       config.RTC
}

func Start(ctx context.Context, config *Config, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgoCli, err := mongoutil.NewMongoDB(ctx, config.Mongo.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.Redis.Build())
	if err != nil {
		return err
	}

	meetingDB, err := mgo2.NewMeetingMongo(mgoCli.GetDB())
	if err != nil {
		return err
	}
	meetingCache := redis.NewMeeting(rdb, meetingDB, redis.GetDefaultOpt())
	database := controller.NewMeeting(meetingDB, meetingCache, mgoCli.GetTx())
	meetingRtc := livekit.NewLiveKit(&config.Rtc)

	u := &meetingServer{
		meetingStorageHandler: database,
		RegisterCenter:        client,
		config:                config,
		meetingRtc:            meetingRtc,
	}
	pbmeeting.RegisterMeetingServiceServer(server, u)
	return nil
}

// BookMeeting Implement the MeetingServiceServer interface
func (s *meetingServer) BookMeeting(ctx context.Context, req *pbmeeting.BookMeetingReq) (*pbmeeting.BookMeetingResp, error) {
	resp := &pbmeeting.BookMeetingResp{}
	meetingDBInfo := &model.MeetingInfo{
		MeetingID:       idutil.OperationIDGenerator(),
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.Scheduled,
		CreatorUserID:   req.CreatorUserID,
	}

	err := s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}
	// fill in response data
	resp.Detail = s.generateRespSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	return resp, nil
}

func (s *meetingServer) CreateImmediateMeeting(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*pbmeeting.CreateImmediateMeetingResp, error) {
	resp := &pbmeeting.CreateImmediateMeetingResp{}
	meetingDBInfo := &model.MeetingInfo{
		MeetingID:       idutil.OperationIDGenerator(),
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.InProgress,
		CreatorUserID:   req.CreatorUserID,
	}
	_, token, liveUrl, err := s.meetingRtc.CreateRoom(ctx, meetingDBInfo.MeetingID)
	if err != nil {
		return resp, err
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}

	metaData := &pbmeeting.MeetingMetadata{}
	meetingDetail := s.generateRespSetting(req.Setting, req.CreatorDefinedMeetingInfo, meetingDBInfo)
	metaData.Detail = meetingDetail
	// create meeting meta data
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	resp.Detail = meetingDetail
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) JoinMeeting(ctx context.Context, req *pbmeeting.JoinMeetingReq) (*pbmeeting.JoinMeetingResp, error) {
	resp := &pbmeeting.JoinMeetingResp{}

	_, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.MeetingID)
	if err != nil {
		return resp, err
	}

	// todo update meta data to liveKit
	//if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
	//	return resp, err
	//}

	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) LeaveMeeting(ctx context.Context, req *pbmeeting.LeaveMeetingReq) (*pbmeeting.LeaveMeetingResp, error) {
	resp := &pbmeeting.LeaveMeetingResp{}

	if err := s.meetingRtc.RemoveParticipant(ctx, req.MeetingID, req.UserID); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *meetingServer) EndMeeting(ctx context.Context, req *pbmeeting.EndMeetingReq) (*pbmeeting.EndMeetingResp, error) {
	resp := &pbmeeting.EndMeetingResp{}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}

	if !s.checkAuthPermission(metaData.Detail.Info.SystemGenerated.CreatorUserID, req.UserID) {
		return resp, errs.ErrArgs.WrapMsg("user did not have permission to end somebody's meeting")
	}

	// change status to completed
	updateData := map[string]any{
		"status": constant.Completed,
	}

	if err := s.meetingRtc.CloseRoom(ctx, req.MeetingID); err != nil {
		return resp, err
	}

	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, updateData); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *meetingServer) GetMeetings(ctx context.Context, req *pbmeeting.GetMeetingsReq) (*pbmeeting.GetMeetingsResp, error) {
	resp := &pbmeeting.GetMeetingsResp{}

	meetings, err := s.meetingStorageHandler.FindByStatus(ctx, req.Status)
	if err != nil {
		return resp, err
	}

	// Create response
	var meetingDetails []*pbmeeting.MeetingInfoSetting
	for _, meeting := range meetings {
		detailSetting, err := s.getMeetingDetailSetting(ctx, meeting)
		if err != nil {
			return resp, err
		}
		meetingDetails = append(meetingDetails, detailSetting)
	}
	resp.MeetingDetails = meetingDetails
	return resp, nil
}

func (s *meetingServer) GetMeeting(ctx context.Context, req *pbmeeting.GetMeetingReq) (*pbmeeting.GetMeetingResp, error) {
	resp := &pbmeeting.GetMeetingResp{}
	meetingDBInfo, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	detailSetting, err := s.getMeetingDetailSetting(ctx, meetingDBInfo)
	if err != nil {
		return resp, err
	}
	resp.MeetingDetail = detailSetting
	return resp, nil
}

func (s *meetingServer) UpdateMeeting(ctx context.Context, req *pbmeeting.UpdateMeetingRequest) (*pbmeeting.UpdateMeetingResp, error) {
	resp := &pbmeeting.UpdateMeetingResp{}

	_, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	// Update the specific field based on the request
	updateData := map[string]any{}
	switch field := req.UpdateField.(type) {
	case *pbmeeting.UpdateMeetingRequest_Title:
		updateData["Title"] = field.Title
	case *pbmeeting.UpdateMeetingRequest_ScheduledTime:
		updateData["ScheduledTime"] = field.ScheduledTime
	case *pbmeeting.UpdateMeetingRequest_MeetingDuration:
		updateData["MeetingDuration"] = field.MeetingDuration
	case *pbmeeting.UpdateMeetingRequest_Password:
		updateData["Password"] = field.Password
	case *pbmeeting.UpdateMeetingRequest_CanParticipantsEnableCamera:
		metaData.Detail.Setting.CanParticipantsEnableCamera = field.CanParticipantsEnableCamera
	case *pbmeeting.UpdateMeetingRequest_CanParticipantsUnmuteMicrophone:
		metaData.Detail.Setting.CanParticipantsUnmuteMicrophone = field.CanParticipantsUnmuteMicrophone
	case *pbmeeting.UpdateMeetingRequest_CanParticipantsShareScreen:
		metaData.Detail.Setting.CanParticipantsShareScreen = field.CanParticipantsShareScreen
	case *pbmeeting.UpdateMeetingRequest_DisableCameraOnJoin:
		metaData.Detail.Setting.DisableCameraOnJoin = field.DisableCameraOnJoin
	case *pbmeeting.UpdateMeetingRequest_DisableMicrophoneOnJoin:
		metaData.Detail.Setting.DisableMicrophoneOnJoin = field.DisableMicrophoneOnJoin
	default:
		return resp, errs.ErrArgs.WrapMsg("unsupported update field")
	}

	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, updateData); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *meetingServer) GetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.GetPersonalMeetingSettingsReq) (*pbmeeting.GetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.GetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}
	for _, personalData := range metaData.PersonalData {
		if personalData.GetUserID() == req.UserID {
			resp.Setting = personalData.PersonalSetting
			break
		}
	}

	return resp, nil
}

func (s *meetingServer) SetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.SetPersonalMeetingSettingsReq) (*pbmeeting.SetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.SetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}
	found := false
	for _, personalData := range metaData.PersonalData {
		if personalData.GetUserID() == req.UserID {
			personalData.PersonalSetting = req.Setting
			found = true
			break
		}
	}
	personalData := &pbmeeting.PersonalData{
		UserID:          req.UserID,
		PersonalSetting: req.Setting,
	}
	if !found {
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}

	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, errs.WrapMsg(err, "update meta data failed")
	}

	return resp, nil
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
