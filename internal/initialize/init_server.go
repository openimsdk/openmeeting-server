package initialize

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"openmeeting-server/internal/initialize/base"
	health_rpc "openmeeting-server/internal/usecase"
	"openmeeting-server/internal/usecase/meeting_rpc"
	"openmeeting-server/protocol/pb"
)

type rtcServer struct {
	meetingSvr meeting_rpc.MeetingGrpc
}

func InitServer(server *grpc.Server) error {
	initFunc := []func() error{
		base.InitLogger,
		base.InitMongo,
		base.InitRedis,
	}

	for _, fc := range initFunc {
		err := fc()
		if err != nil {
			return err
		}
	}

	// register
	context := context.Background()
	meetingGrpc := meeting_rpc.NewMeetingGrpc(context)
	if meetingGrpc == nil {
		return errors.New("init meeting grpc failed")
	}
	healthGrpc := health_rpc.NewHealthCheckGrpc()
	pb.RegisterHealthServer(server, healthGrpc)
	pb.RegisterMeetingServiceServer(server, meetingGrpc)
	fmt.Println("start rtc service success!")
	return nil
}
