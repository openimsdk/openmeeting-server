package health_rpc

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"openmeeting-server/protocol/pb"
)

type HealthCheckGrpc struct {
}

func NewHealthCheckGrpc() *HealthCheckGrpc {
	return &HealthCheckGrpc{}
}

func (s *HealthCheckGrpc) Ping(context.Context, *pb.PingBody) (*pb.PingBody, error) {
	reply := "pong"
	return &pb.PingBody{
		Str: reply,
	}, errs.ErrRecordNotFound
}
