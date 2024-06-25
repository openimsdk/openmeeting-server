package rpcclient

import (
	"context"
	"github.com/openimsdk/protocol/openmeeting/meeting"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type Meeting struct {
	conn      grpc.ClientConnInterface
	Client    meeting.MeetingServiceClient
	Discovery discovery.SvcDiscoveryRegistry
}

// NewMeeting initializes and returns a User instance based on the provided service discovery registry.
func NewMeeting(discovery discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Meeting {
	conn, err := discovery.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := meeting.NewMeetingServiceClient(conn)
	return &Meeting{Discovery: discovery, Client: client,
		conn: conn,
	}
}
