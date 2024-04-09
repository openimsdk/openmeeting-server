package meeting_rpc

import (
	"context"
	meeting_domain "openmeeting-server/internal/domain/meeting"
	"openmeeting-server/protocol/pb"
)

type MeetingGrpc struct {
	meetingDomain meeting_domain.MeetingDomainInterface
}

func NewMeetingGrpc() *MeetingGrpc {
	return &MeetingGrpc{
		meetingDomain: meeting_domain.NewMeetingService(),
	}
}

func (mg *MeetingGrpc) DeleteMeetingRecords(ctx context.Context, request *pb.DeleteMeetingRecordsReq) (*pb.DeleteMeetingRecordsResp, error) {
	err := mg.meetingDomain.DeleteMeeting(ctx, request.RoomIDs...)
	if err != nil {
		return &pb.DeleteMeetingRecordsResp{}, err
	}
	return &pb.DeleteMeetingRecordsResp{}, nil
}
