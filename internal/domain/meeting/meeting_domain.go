package meeting_domain

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/tx"
	"openmeeting-server/internal/infrastructure/cache"
	"openmeeting-server/internal/infrastructure/repository"
)

type MeetingDomainInterface interface {
	CreateMeeting(ctx context.Context) error
	DeleteMeeting(ctx context.Context, roomId ...string) error
}

type MeetingDomain struct {
	MeetingRepository repository.MeetingInterface

	tx tx.CtxTx
}

func NewMeetingService() *MeetingDomain {
	repo, err := repository.NewMeetingRepository()
	if err != nil {
		return nil
	}
	client, err := cache.GetMongoClient()
	if err != nil {
		return nil
	}

	return &MeetingDomain{
		MeetingRepository: repo,
		tx:                tx.NewMongo(client),
	}
}

func (m *MeetingDomain) CreateMeeting(ctx context.Context) error {
	return nil
}

func (m *MeetingDomain) DeleteMeeting(ctx context.Context, roomId ...string) error {
	err := m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.MeetingRepository.DeleteMeetingInfos(ctx, roomId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errs.ErrDatabase
	}
	return nil
}
