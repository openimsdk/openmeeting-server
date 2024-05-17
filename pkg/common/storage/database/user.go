package database

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
)

type User interface {
	Create(ctx context.Context, users []*model.User) (err error)
	Take(ctx context.Context, userID string) (user *model.User, err error)
}
