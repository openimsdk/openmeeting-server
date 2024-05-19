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
	// TakeWithError Get the information of the specified meeting. If the meetingID is not found, it will also return an error
	TakeWithError(ctx context.Context, meetingID string) (meeting *model.MeetingInfo, err error) //1
	// Create Insert multiple external guarantees that the meetingID is not repeated and does not exist in the storage
	Create(ctx context.Context, meetings []*model.MeetingInfo) (err error) //1
	Update(ctx context.Context, meetingID string, updateData map[string]any) (err error)
	FindByStatus(ctx context.Context, status string) ([]*model.MeetingInfo, error)
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
func (u *MeetingStorageManager) TakeWithError(ctx context.Context, meetingID string) (meeting *model.MeetingInfo, err error) {
	meeting, err = u.db.Take(ctx, meetingID)
	if err != nil {
		return meeting, errs.WrapMsg(err, "get record from mongo failed, meetingID: ", meetingID)
	}
	err = errs.ErrRecordNotFound.WrapMsg("meetingID not found")
	return
}

// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the storage.
func (u *MeetingStorageManager) Create(ctx context.Context, meetings []*model.MeetingInfo) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.db.Create(ctx, meetings); err != nil {
			return errs.WrapMsg(err, "create meeting data failed")
		}
		return nil
	})
}

func (u *MeetingStorageManager) Update(ctx context.Context, meetingID string, updateData map[string]any) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.db.Update(ctx, meetingID, updateData); err != nil {
			return errs.WrapMsg(err, "update meeting info failed, meetingID:", meetingID)
		}
		return nil
	})
}

func (u *MeetingStorageManager) FindByStatus(ctx context.Context, status string) ([]*model.MeetingInfo, error) {
	return u.db.FindByStatus(ctx, status)
}
