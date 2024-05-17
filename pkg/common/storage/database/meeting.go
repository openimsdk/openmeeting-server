package database

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
)

type Meeting interface {
	Create(ctx context.Context, meetings []*model.MeetingInfo) (err error)
	Take(ctx context.Context, meetingID string) (meeting *model.MeetingInfo, err error)
}
