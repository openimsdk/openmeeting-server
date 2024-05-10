package meeting

import (
	"context"
	"google.golang.org/grpc"
	meeting_domain "openmeeting-server/internal/domain/meeting"
	"openmeeting-server/pkg/common/config"
	"openmeeting-server/protocol/pb"
)

type meetingServer struct {
	meetingDomain meeting_domain.MeetingDomainInterface
}

func newMeetingServer(ctx context.Context, serverConfig *config.Config) (*meetingServer, error) {
	meetingDomain, err := meeting_domain.NewMeetingService(ctx, serverConfig)
	if err != nil {
		return nil, err
	}
	return &meetingServer{
		meetingDomain: meetingDomain,
	}, nil
}

func Start(ctx context.Context, serverConfig config.Config, server *grpc.Server) error {
	mServer, err := newMeetingServer(ctx, &serverConfig)
	if err != nil {
		return err
	}
	pb.RegisterMeetingServiceServer(server, mServer)
	return nil
}

func (mg *meetingServer) GetMeetingList(ctx context.Context, request *pb.GetMeetingListReq) (*pb.GetMeetingListResp, error) {
	return nil, nil
}

func (mg *meetingServer) GetMeetingDetailInfo(ctx context.Context, request *pb.GetMeetingDetailInfoReq) (*pb.GetMeetingDetailInfoResp, error) {
	return nil, nil
}

func (mg *meetingServer) CreateQuickMeeting(ctx context.Context, request *pb.QuickCreateMeetingReq) (*pb.QuickCreateMeetingResp, error) {
	resp, err := mg.meetingDomain.QuickStartCreateMeeting(ctx, request)
	if err != nil {
		return &pb.QuickCreateMeetingResp{}, err
	}
	return resp, err
}

func (mg *meetingServer) PreBookCreateMeeting(ctx context.Context, request *pb.PreBookCreateMeetingReq) (*pb.PreBookCreateMeetingResp, error) {
	resp, err := mg.meetingDomain.PreBookCreateMeeting(ctx, request)
	if err != nil {
		return &pb.PreBookCreateMeetingResp{}, err
	}
	return resp, nil
}

func (mg *meetingServer) UpdatePreBookMeeting(ctx context.Context, request *pb.PreBookUpdateMeetingReq) (*pb.PreBookUpdateMeetingResp, error) {
	return mg.meetingDomain.UpdatePreBookMeeting(ctx, request)
}

func (mg *meetingServer) JoinMeeting(ctx context.Context, request *pb.JoinMeetingReq) (*pb.JoinMeetingResp, error) {
	resp, err := mg.meetingDomain.JoinMeeting(ctx, request)
	if err != nil {
		return &pb.JoinMeetingResp{
			MeetingID: request.MeetingID,
		}, err
	}
	return resp, nil
}

func (mg *meetingServer) DeleteMeeting(ctx context.Context, request *pb.DeleteMeetingReq) (*pb.DeleteMeetingResp, error) {
	err := mg.meetingDomain.DeleteMeeting(ctx, request.MeetingID)
	if err != nil {
		return &pb.DeleteMeetingResp{}, err
	}
	return &pb.DeleteMeetingResp{}, nil
}

func (mg *meetingServer) UpdateMeetingInfo(ctx context.Context, request *pb.UpdateMeetingInfoReq) (*pb.UpdateMeetingInfoResp, error) {
	return mg.meetingDomain.UpdateMeetingInfo(ctx, request)
}

func (mg *meetingServer) ToggleMeetingMedia(ctx context.Context, request *pb.ToggleMeetingMediaReq) (*pb.ToggleMeetingMediaResp, error) {
	return nil, nil
}

func (mg *meetingServer) ManageMeetingUserMedia(ctx context.Context, request *pb.ManageMeetingMediaReq) (*pb.ManageMeetingMediaResp, error) {
	return nil, nil
}

func (mg *meetingServer) UpdateMeetingAction(context.Context, *pb.UpdateMeetingActionReq) (*pb.UpdateMeetingActionResp, error) {
	return nil, nil
}

func (mg *meetingServer) CloseMeeting(ctx context.Context, request *pb.CloseMeetingReq) (*pb.CloseMeetingResp, error) {
	resp, err := mg.meetingDomain.CloseMeeting(ctx, request)
	if err != nil {
		return &pb.CloseMeetingResp{}, err
	}
	return resp, nil
}

func (mg *meetingServer) LeaveMeeting(ctx context.Context, request *pb.LeaveMeetingReq) (*pb.LeaveMeetingResp, error) {
	resp, err := mg.meetingDomain.LeaveMeeting(ctx, request)
	if err != nil {
		return &pb.LeaveMeetingResp{}, err
	}
	return resp, nil
}

func (mg *meetingServer) KickOffMeeting(ctx context.Context, request *pb.KickOffMeetingReq) (*pb.KickOffMeetingResp, error) {
	resp, err := mg.meetingDomain.KickOffMeeting(ctx, request)
	if err != nil {
		return &pb.KickOffMeetingResp{}, err
	}
	return resp, nil
}
