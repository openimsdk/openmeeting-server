package meeting_rpc

import (
	"context"
	meeting_domain "openmeeting-server/internal/domain/meeting"
	"openmeeting-server/protocol/pb"
)

type MeetingGrpc struct {
	meetingDomain meeting_domain.MeetingDomainInterface
}

func NewMeetingGrpc(ctx context.Context) *MeetingGrpc {
	meetingDomain, err := meeting_domain.NewMeetingService(ctx)
	if err != nil {
		return nil
	}
	return &MeetingGrpc{
		meetingDomain: meetingDomain,
	}
}

func (mg *MeetingGrpc) GetMeetingList(ctx context.Context, request *pb.GetMeetingListReq) (*pb.GetMeetingListResp, error) {
	return nil, nil
}

func (mg *MeetingGrpc) GetMeetingDetailInfo(ctx context.Context, request *pb.GetMeetingDetailInfoReq) (*pb.GetMeetingDetailInfoResp, error) {
	return nil, nil
}

func (mg *MeetingGrpc) CreateQuickMeeting(ctx context.Context, request *pb.QuickCreateMeetingReq) (*pb.QuickCreateMeetingResp, error) {
	resp, err := mg.meetingDomain.QuickStartCreateMeeting(ctx, request)
	if err != nil {
		return &pb.QuickCreateMeetingResp{}, err
	}
	return resp, err
}

func (mg *MeetingGrpc) PreBookCreateMeeting(ctx context.Context, request *pb.PreBookCreateMeetingReq) (*pb.PreBookCreateMeetingResp, error) {
	resp, err := mg.meetingDomain.PreBookCreateMeeting(ctx, request)
	if err != nil {
		return &pb.PreBookCreateMeetingResp{}, err
	}
	return resp, nil
}

func (mg *MeetingGrpc) UpdatePreBookMeeting(ctx context.Context, request *pb.PreBookUpdateMeetingReq) (*pb.PreBookUpdateMeetingResp, error) {
	return mg.meetingDomain.UpdatePreBookMeeting(ctx, request)
}

func (mg *MeetingGrpc) JoinMeeting(ctx context.Context, request *pb.JoinMeetingReq) (*pb.JoinMeetingResp, error) {
	resp, err := mg.meetingDomain.JoinMeeting(ctx, request)
	if err != nil {
		return &pb.JoinMeetingResp{
			MeetingID: request.MeetingID,
		}, err
	}
	return resp, nil
}

func (mg *MeetingGrpc) DeleteMeeting(ctx context.Context, request *pb.DeleteMeetingReq) (*pb.DeleteMeetingResp, error) {
	err := mg.meetingDomain.DeleteMeeting(ctx, request.MeetingID)
	if err != nil {
		return &pb.DeleteMeetingResp{}, err
	}
	return &pb.DeleteMeetingResp{}, nil
}

func (mg *MeetingGrpc) UpdateMeetingInfo(ctx context.Context, request *pb.UpdateMeetingInfoReq) (*pb.UpdateMeetingInfoResp, error) {
	return mg.meetingDomain.UpdateMeetingInfo(ctx, request)
}

func (mg *MeetingGrpc) ToggleMeetingMedia(ctx context.Context, request *pb.ToggleMeetingMediaReq) (*pb.ToggleMeetingMediaResp, error) {
	return nil, nil
}

func (mg *MeetingGrpc) ManageMeetingUserMedia(ctx context.Context, request *pb.ManageMeetingMediaReq) (*pb.ManageMeetingMediaResp, error) {
	return nil, nil
}

func (mg *MeetingGrpc) UpdateMeetingAction(context.Context, *pb.UpdateMeetingActionReq) (*pb.UpdateMeetingActionResp, error) {
	return nil, nil
}

func (mg *MeetingGrpc) CloseMeeting(ctx context.Context, request *pb.CloseMeetingReq) (*pb.CloseMeetingResp, error) {
	resp, err := mg.meetingDomain.CloseMeeting(ctx, request)
	if err != nil {
		return &pb.CloseMeetingResp{}, err
	}
	return resp, nil
}

func (mg *MeetingGrpc) LeaveMeeting(ctx context.Context, request *pb.LeaveMeetingReq) (*pb.LeaveMeetingResp, error) {
	resp, err := mg.meetingDomain.LeaveMeeting(ctx, request)
	if err != nil {
		return &pb.LeaveMeetingResp{}, err
	}
	return resp, nil
}

func (mg *MeetingGrpc) KickOffMeeting(ctx context.Context, request *pb.KickOffMeetingReq) (*pb.KickOffMeetingResp, error) {
	resp, err := mg.meetingDomain.KickOffMeeting(ctx, request)
	if err != nil {
		return &pb.KickOffMeetingResp{}, err
	}
	return resp, nil
}
