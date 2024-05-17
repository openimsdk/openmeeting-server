package controller

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/errs"
)

type Meeting interface {
	// TakeWithError Get the information of the specified user. If the userID is not found, it will also return an error
	TakeWithError(ctx context.Context, meetingID string) (meeting *model.MeetingInfo, err error) //1
	// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the storage
	Create(ctx context.Context, meetings []*model.MeetingInfo) (err error) //1

}

type MeetingStorageManager struct {
	tx    tx.Tx
	db    database.Meeting
	cache cache.Meeting
}

func NewMeeting(meetingDB database.Meeting, cache cache.Meeting, tx tx.Tx) Meeting {
	return &MeetingStorageManager{db: meetingDB, cache: cache, tx: tx}
}

// TakeWithError Get the information of the specified user and return an error if the userID is not found.
func (u *MeetingStorageManager) TakeWithError(ctx context.Context, meetingsID string) (meeting *model.MeetingInfo, err error) {
	meeting, err = u.db.Take(ctx, meetingsID)
	if err != nil {
		return
	}
	err = errs.ErrRecordNotFound.WrapMsg("userID not found")
	return
}

// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the storage.
func (u *MeetingStorageManager) Create(ctx context.Context, meetings []*model.MeetingInfo) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.db.Create(ctx, meetings); err != nil {
			return err
		}
		return nil
	})
}
