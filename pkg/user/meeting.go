package user

import (
	"context"
	"github.com/openimsdk/protocol/openmeeting/user"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
)

func NewMeetingUserClient(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) user.UserClient {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	return user.NewUserClient(conn)
}

func NewMeeting(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) User {
	return &meeting{user: NewMeetingUserClient(discov, rpcRegisterName)}
}

type meeting struct {
	user user.UserClient
}

func (m *meeting) GetUsersInfos(ctx context.Context, userIDs []string) ([]*user.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*user.UserInfo{}, nil
	}
	resp, err := m.user.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.UsersInfo, nil
}
