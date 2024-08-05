package user

import (
	"context"
	"github.com/openimsdk/protocol/openmeeting/user"
)

type User interface {
	GetUsersInfos(ctx context.Context, userIDs []string) ([]*user.UserInfo, error)
}
