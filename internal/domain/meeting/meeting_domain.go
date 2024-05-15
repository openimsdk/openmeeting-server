package meeting_domain

import (
	"context"
	"errors"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/timeutil"
	"openmeeting-server/internal/domain/livekit"
	"openmeeting-server/internal/infrastructure/repository"
	"openmeeting-server/internal/infrastructure/repository/model"
	"openmeeting-server/internal/utils"
	"openmeeting-server/pkg/common/config"
	"openmeeting-server/protocol/pb"
	"time"
)

type MeetingDomainInterface interface {
	PreBookCreateMeeting(ctx context.Context, req *pb.PreBookCreateMeetingReq) (*pb.PreBookCreateMeetingResp, error)
	UpdatePreBookMeeting(ctx context.Context, req *pb.PreBookUpdateMeetingReq) (*pb.PreBookUpdateMeetingResp, error)
	DeleteMeeting(ctx context.Context, roomId ...string) error
	QuickStartCreateMeeting(ctx context.Context, req *pb.QuickCreateMeetingReq) (*pb.QuickCreateMeetingResp, error)
	JoinMeeting(ctx context.Context, req *pb.JoinMeetingReq) (*pb.JoinMeetingResp, error)

	UpdateMeetingInfo(ctx context.Context, req *pb.UpdateMeetingInfoReq) (*pb.UpdateMeetingInfoResp, error)

	CloseMeeting(ctx context.Context, req *pb.CloseMeetingReq) (*pb.CloseMeetingResp, error)
	LeaveMeeting(ctx context.Context, req *pb.LeaveMeetingReq) (*pb.LeaveMeetingResp, error)
	KickOffMeeting(ctx context.Context, req *pb.KickOffMeetingReq) (*pb.KickOffMeetingResp, error)
}

type MeetingDomain struct {
	MeetingRepository repository.MeetingInterface
	RTCLogic          livekit.RTCDomain
	serverConfig      *config.Config
	mongoClient       *mongoutil.Client
}

func NewMeetingService(ctx context.Context, c *config.Config) (*MeetingDomain, error) {
	mongoClient, err := mongoutil.NewMongoDB(ctx, c.MongodbConfig.Build())
	if err != nil {
		return nil, err
	}

	repo, err := repository.NewMeetingRepository(mongoClient.GetDB())
	if err != nil {
		return nil, err
	}

	return &MeetingDomain{
		MeetingRepository: repo,
		RTCLogic:          livekit.NewLiveKit(ctx, c),
		serverConfig:      c,
		mongoClient:       mongoClient,
	}, nil
}

func (m *MeetingDomain) PreBookCreateMeeting(ctx context.Context, req *pb.PreBookCreateMeetingReq) (*pb.PreBookCreateMeetingResp, error) {

	meetingID := utils.GenerateUniqueKey()

	info := &model.MeetingInfo{
		MeetingID:     meetingID,
		MeetingName:   req.MeetingName,
		HostUserID:    req.UserID,
		CreatorUserID: req.UserID,
		StartTime:     timeutil.UnixSecondToTime(req.StartTime),
		EndTime:       timeutil.UnixSecondToTime(req.StartTime + req.MeetingDuration),
		Duration:      req.MeetingDuration,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := m.MeetingRepository.CreateMeetingInfo(ctx, info); err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *MeetingDomain) UpdatePreBookMeeting(ctx context.Context, req *pb.PreBookUpdateMeetingReq) (*pb.PreBookUpdateMeetingResp, error) {
	updateData := map[string]any{"status": false}
	if err := m.MeetingRepository.UpdateMeetingInfo(ctx, req.MeetingId, updateData); err != nil {

	}

	return nil, nil
}

func (m *MeetingDomain) UpdateMeetingInfo(ctx context.Context, req *pb.UpdateMeetingInfoReq) (*pb.UpdateMeetingInfoResp, error) {
	info := &pb.MeetingInfo{
		RoomID: req.MeetingID,
	}

	// update meta data
	if err := m.RTCLogic.UpdateMetaData(ctx, info); err != nil {
		return &pb.UpdateMeetingInfoResp{}, err
	}
	// update mongo db

	return &pb.UpdateMeetingInfoResp{}, nil
}

func (m *MeetingDomain) DeleteMeeting(ctx context.Context, roomId ...string) error {
	err := m.mongoClient.GetTx().Transaction(ctx, func(ctx context.Context) error {
		if err := m.MeetingRepository.DeleteMeetingInfos(ctx, roomId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errs.ErrInternalServer.WrapMsg("delete meeting failed: ", roomId)
	}
	return nil
}

func (m *MeetingDomain) QuickStartCreateMeeting(ctx context.Context, req *pb.QuickCreateMeetingReq) (*pb.QuickCreateMeetingResp, error) {
	//if _, err := x.roomIsExist(ctx, req.RoomID); err != nil && errs.Unwrap(err) != errs.ErrRecordNotFound {
	//	return nil, err
	//}
	//
	meetingID := utils.GenerateUniqueKey()
	_, token, liveUrl, err := m.RTCLogic.CreateRoom(ctx, meetingID)
	if err != nil {
		return nil, errs.WrapMsg(err, "create room failed")
	}

	info := &model.MeetingInfo{
		MeetingID:     meetingID,
		MeetingName:   req.MeetingName,
		HostUserID:    req.UserID,
		CreatorUserID: req.UserID,
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		Duration:      600, // default configuration
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := m.MeetingRepository.CreateMeetingInfo(ctx, info); err != nil {
		return nil, err
	}

	return &pb.QuickCreateMeetingResp{
		MeetingID: meetingID,
		Token:     token,
		LiveURL:   liveUrl,
	}, nil
}

func (m *MeetingDomain) JoinMeeting(ctx context.Context, req *pb.JoinMeetingReq) (*pb.JoinMeetingResp, error) {
	metaData, err := m.RTCLogic.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, errs.WrapMsg(err, "get room data failed")
	}

	log.ZDebug(ctx, "metaData", metaData)

	token, liveUrl, err := m.RTCLogic.GetJoinToken(ctx, req.MeetingID, req.MeetingID)
	if err != nil {
		return nil, errs.WrapMsg(err, "get join token failed")
	}

	if err := m.RTCLogic.UpdateMetaData(ctx, &pb.MeetingInfo{
		RoomID: req.MeetingID,
	}); err != nil {
		return nil, errs.WrapMsg(err, "update meta data failed")
	}

	return &pb.JoinMeetingResp{
		MeetingID: req.MeetingID,
		Token:     token,
		LiveURL:   liveUrl,
	}, nil
}

func (m *MeetingDomain) CloseMeeting(ctx context.Context, req *pb.CloseMeetingReq) (*pb.CloseMeetingResp, error) {
	metaData, err := m.RTCLogic.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}

	if !m.checkAuthPermission(metaData.HostUserID, req.UserID) {
		return nil, errors.New("user did not have permission to close meeting")
	}

	if err := m.RTCLogic.CloseRoom(ctx, req.MeetingID); err != nil {
		return nil, err
	}

	return &pb.CloseMeetingResp{}, nil
}

func (m *MeetingDomain) checkAuthPermission(hostUserID, requestUserID string) bool {
	return hostUserID == requestUserID
}

func (m *MeetingDomain) LeaveMeeting(ctx context.Context, req *pb.LeaveMeetingReq) (*pb.LeaveMeetingResp, error) {

	if err := m.RTCLogic.RemoveParticipant(ctx, req.MeetingID, req.LeaveUserID); err != nil {
		return &pb.LeaveMeetingResp{}, nil
	}
	return &pb.LeaveMeetingResp{}, nil
}

func (m *MeetingDomain) KickOffMeeting(ctx context.Context, req *pb.KickOffMeetingReq) (*pb.KickOffMeetingResp, error) {
	metaData, err := m.RTCLogic.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}

	if !m.checkAuthPermission(metaData.HostUserID, req.KickerUserID) {
		return nil, errors.New("user did not have permission to kick somebody out of the meeting")
	}

	if err := m.RTCLogic.RemoveParticipant(ctx, req.MeetingID, req.LeaveUserID); err != nil {
		return &pb.KickOffMeetingResp{}, nil
	}
	return &pb.KickOffMeetingResp{}, nil
}
