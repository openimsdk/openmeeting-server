package meeting

import (
	"context"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	mgo2 "github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	pbmeeting "github.com/openimsdk/openmeeting-server/pkg/protocol/meeting"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	registry "github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/utils/idutil"
	"google.golang.org/grpc"
)

type meetingServer struct {
	meetingStorageHandler controller.Meeting
	RegisterCenter        registry.SvcDiscoveryRegistry
	config                *Config
}

type Config struct {
	Rpc       config.User
	Redis     config.Redis
	Mongo     config.Mongo
	Discovery config.Discovery
	Share     config.Share
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
	u := &meetingServer{
		meetingStorageHandler: database,
		RegisterCenter:        client,
		config:                config,
	}
	pbmeeting.RegisterMeetingServiceServer(server, u)
	return nil
}

// BookMeeting Implement the MeetingServiceServer interface
func (s *meetingServer) BookMeeting(ctx context.Context, req *pbmeeting.BookMeetingReq) (*pbmeeting.BookMeetingResp, error) {
	fmt.Println("BookMeeting called")
	resp := &pbmeeting.BookMeetingResp{}
	meetingDBInfo := &model.MeetingInfo{
		MeetingID:       idutil.OperationIDGenerator(),
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		CreatorUserID:   req.CreatorUserID,
	}

	err := s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}
	// fill in response data
	resp.Detail = s.generateRespSetting(req, meetingDBInfo)
	return resp, nil
}

// generateMeetingInfoSetting generates MeetingInfoSetting from the given request and meeting ID.
func (s *meetingServer) generateRespSetting(req *pbmeeting.BookMeetingReq, meeting *model.MeetingInfo) *pbmeeting.MeetingInfoSetting {
	// Fill in response data
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID: req.CreatorUserID,
		Status:        meeting.Status,
		StartTime:     req.CreatorDefinedMeetingInfo.ScheduledTime, // Scheduled start time as the actual start time
		MeetingID:     meeting.MeetingID,
	}
	// Combine system-generated and creator-defined info
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: req.CreatorDefinedMeetingInfo,
	}
	// Create MeetingInfoSetting
	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info:    meetingInfo,
		Setting: req.Setting,
	}
	return meetingInfoSetting
}

func (s *meetingServer) CreateImmediateMeeting(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*pbmeeting.CreateImmediateMeetingResp, error) {
	fmt.Println("CreateImmediateMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.CreateImmediateMeetingResp{}, nil
}

func (s *meetingServer) JoinMeeting(ctx context.Context, req *pbmeeting.JoinMeetingReq) (*pbmeeting.JoinMeetingResp, error) {
	fmt.Println("JoinMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.JoinMeetingResp{}, nil
}

func (s *meetingServer) LeaveMeeting(ctx context.Context, req *pbmeeting.LeaveMeetingReq) (*pbmeeting.LeaveMeetingResp, error) {
	fmt.Println("LeaveMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.LeaveMeetingResp{}, nil
}

func (s *meetingServer) EndMeeting(ctx context.Context, req *pbmeeting.EndMeetingReq) (*pbmeeting.EndMeetingResp, error) {
	fmt.Println("EndMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.EndMeetingResp{}, nil
}

func (s *meetingServer) GetMeetings(ctx context.Context, req *pbmeeting.GetMeetingsReq) (*pbmeeting.GetMeetingsResp, error) {
	fmt.Println("GetMeetings called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.GetMeetingsResp{}, nil
}

func (s *meetingServer) GetMeeting(ctx context.Context, req *pbmeeting.GetMeetingReq) (*pbmeeting.GetMeetingResp, error) {
	fmt.Println("GetMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.GetMeetingResp{}, nil
}

func (s *meetingServer) UpdateMeeting(ctx context.Context, req *pbmeeting.UpdateMeetingRequest) (*pbmeeting.UpdateMeetingResp, error) {
	fmt.Println("UpdateMeeting called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.UpdateMeetingResp{}, nil
}

func (s *meetingServer) GetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.GetPersonalMeetingSettingsReq) (*pbmeeting.GetPersonalMeetingSettingsResp, error) {
	fmt.Println("GetPersonalMeetingSettings called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.GetPersonalMeetingSettingsResp{}, nil
}

func (s *meetingServer) SetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.SetPersonalMeetingSettingsReq) (*pbmeeting.SetPersonalMeetingSettingsResp, error) {
	fmt.Println("SetPersonalMeetingSettings called")
	// Add logic to handle the request using s.userStorageHandler, s.RegisterCenter, and s.config
	return &pbmeeting.SetPersonalMeetingSettingsResp{}, nil
}
