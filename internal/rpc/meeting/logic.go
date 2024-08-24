package meeting

import (
	"context"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common/constant"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/openimsdk/tools/utils/timeutil"
	"time"
)

const (
	video = "video"
	audio = "audio"
)

func (s *meetingServer) setSelfPersonalSetting(ctx context.Context, metaData *pbmeeting.MeetingMetadata, req *pbmeeting.SetPersonalMeetingSettingsReq) error {
	found := false
	var personalData *pbmeeting.PersonalData
	for _, one := range metaData.PersonalData {
		if one.GetUserID() == req.UserID {
			personalData = one
			found = true
			break
		}
	}
	needUpdate := true
	if found {
		if req.CameraOnEntry == nil && req.MicrophoneOnEntry == nil {
			needUpdate = false
		}
	}
	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
	}
	if req.CameraOnEntry != nil {
		personalData.PersonalSetting.CameraOnEntry = req.CameraOnEntry.Value
	}
	if req.MicrophoneOnEntry != nil {
		personalData.PersonalSetting.MicrophoneOnEntry = req.MicrophoneOnEntry.Value
	}
	toggle := s.checkUserEnableCamera(metaData.Detail.Setting, personalData)
	if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, video, !toggle); err != nil {
		return errs.WrapMsg(err, "toggle camera stream failed")
	}

	toggle = s.checkUserEnableMicrophone(metaData.Detail.Setting, personalData)
	if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, audio, !toggle); err != nil {
		return errs.WrapMsg(err, "toggle microphone stream failed")
	}

	// judge whether user need to change or not
	if !found {
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}
	if !needUpdate {
		log.CInfo(ctx, "no need update meta data for set setting")
		// no need update
		return nil
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return errs.WrapMsg(err, "update meta data failed")
	}

	return nil
}

// setParticipantPersonalSetting set setting
func (s *meetingServer) setParticipantPersonalSetting(ctx context.Context, metaData *pbmeeting.MeetingMetadata, req *pbmeeting.SetPersonalMeetingSettingsReq) error {
	found := false
	var personalData *pbmeeting.PersonalData
	for _, one := range metaData.PersonalData {
		if one.GetUserID() == req.UserID {
			personalData = one
			found = true
			break
		}
	}

	// judge whether user need to change or not
	needUpdate := true
	if found {
		if req.CameraOnEntry == nil && req.MicrophoneOnEntry == nil {
			needUpdate = false
		}
	}

	if !found {
		personalData = s.generateDefaultPersonalData(req.UserID)
	}
	if req.CameraOnEntry != nil {
		personalData.LimitSetting.CameraOnEntry = req.CameraOnEntry.Value
	}
	if req.MicrophoneOnEntry != nil {
		personalData.LimitSetting.MicrophoneOnEntry = req.MicrophoneOnEntry.Value
	}
	var (
		cameraOn     = false
		microphoneOn = false
	)
	if cameraOn = s.checkUserEnableCamera(metaData.Detail.Setting, personalData); !cameraOn {
		// no need to care the scene that turning on the camera
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, video, true); err != nil {
			return errs.WrapMsg(err, "toggle camera stream failed")
		}
	}

	if microphoneOn = s.checkUserEnableMicrophone(metaData.Detail.Setting, personalData); !microphoneOn {
		// no need to care the scene that turning on the microphone
		if err := s.meetingRtc.ToggleMimeStream(ctx, req.MeetingID, req.UserID, audio, true); err != nil {
			return errs.WrapMsg(err, "toggle microphone stream failed")
		}
	}

	if !found {
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}
	if !needUpdate {
		// no need to update
		log.CInfo(ctx, "no need update meta data for set setting")
		return nil
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return errs.WrapMsg(err, "update meta data failed")
	}

	if err := s.sendStreamOperateData2Client(ctx, req.MeetingID, req.UserID, req.CameraOnEntry, req.MicrophoneOnEntry); err != nil {
		return errs.WrapMsg(err, "send data failed")
	}

	return nil
}

func (s *meetingServer) muteAllStream(ctx context.Context, roomID, streamType string, mute bool) (streamNotExistUserIDList []string, failedUserIDList []string, err error) {
	participants, err := s.meetingRtc.ListParticipants(ctx, roomID)
	if err != nil {
		return nil, nil, errs.WrapMsg(err, "get participant list failed")
	}
	for _, v := range participants {
		err := s.meetingRtc.ToggleMimeStream(ctx, roomID, v.Identity, streamType, mute)
		if err != nil {
			log.ZError(ctx, "muteAllStream failed", err)
			if errs.ErrRecordNotFound.Is(err) {
				streamNotExistUserIDList = append(streamNotExistUserIDList, v.Identity)
			} else {
				failedUserIDList = append(failedUserIDList, v.Identity)
			}
		}
	}
	return
}

func (s *meetingServer) refreshNonRepeatMeeting(ctx context.Context, info *model.MeetingInfo) map[string]any {
	updateData := map[string]any{}
	nowTimestamp, err := timeutil.GetTimestampByTimezone(info.TimeZone)
	if err != nil {
		return updateData
	}
	if info.StartTime+info.MeetingDuration < nowTimestamp {
		updateData["status"] = constant.Completed
	} else if info.StartTime+info.MeetingDuration < nowTimestamp && info.StartTime > nowTimestamp {
		if info.Status == constant.InProgress {
			return updateData
		}
		updateData["status"] = constant.InProgress
	}
	return updateData
}

func (s *meetingServer) IsTodayNeedMeeting(ctx context.Context, info *model.MeetingInfo) bool {
	currentTimestamp, err := timeutil.GetTimestampByTimezone(info.TimeZone)
	if err != nil {
		return false
	}
	switch info.RepeatType {
	case constant.RepeatDaily:
		// always return true
		return true
	case constant.RepeatWeekly:
		// judge the day is the same day
		sameDay, err := timeutil.IsSameWeekday(info.TimeZone, info.ScheduledTime)
		if err != nil {
			log.ZError(ctx, "error", err)
			return false
		}
		return sameDay
	case constant.RepeatMonth:
		// judge the day is the same day of one month
		sameDay, err := timeutil.IsSameDayOfMonth(info.TimeZone, info.ScheduledTime)
		if err != nil {
			log.ZError(ctx, "error", err)
			return false
		}
		return sameDay
	case constant.RepeatWeekDay:
		// judge the day is weekday
		return timeutil.IsWeekday(currentTimestamp)

	case constant.RepeatCustom:
		// customize repeat days
		if len(info.RepeatDayOfWeek) > 0 {
			startTime := time.Unix(currentTimestamp, 0)
			if stringutil.IsContainInt32(int32(startTime.Day()), info.RepeatDayOfWeek) {
				return true
			}
			return false
		}
		return processIntervalType(ctx, info)
	default:
		return false
	}
}

func checkCycle(ctx context.Context, functions ...func() (bool, error)) bool {
	for _, fn := range functions {
		result, err := fn()
		if err != nil {
			log.ZError(ctx, "error", err)
			return false
		}
		// return true only when all functions return true
		if !result {
			return false
		}
	}
	return true
}

func checkNthMonthCycle(ctx context.Context, info *model.MeetingInfo) bool {
	return checkCycle(ctx,
		func() (bool, error) {
			return timeutil.IsNthMonthCycle(info.TimeZone, info.ScheduledTime, int(info.Interval))
		},
		func() (bool, error) {
			return timeutil.IsSameDayOfMonth(info.TimeZone, info.ScheduledTime)
		})
}

func checkNthWeekCycle(ctx context.Context, info *model.MeetingInfo) bool {
	return checkCycle(ctx,
		func() (bool, error) {
			return timeutil.IsNthWeekCycle(info.TimeZone, info.ScheduledTime, int(info.Interval))
		},
		func() (bool, error) {
			return timeutil.IsSameWeekday(info.TimeZone, info.ScheduledTime)
		})
}

func checkNthDayCycle(ctx context.Context, info *model.MeetingInfo) bool {
	return checkCycle(ctx,
		func() (bool, error) {
			return timeutil.IsNthDayCycle(info.TimeZone, info.ScheduledTime, int(info.Interval))
		})
}

func processIntervalType(ctx context.Context, info *model.MeetingInfo) bool {
	switch info.UintType {
	case constant.UnitTypeMonth:
		return checkNthMonthCycle(ctx, info)
	case constant.UnitTypeWeek:
		return checkNthWeekCycle(ctx, info)
	case constant.UnitTypeDay:
		return checkNthDayCycle(ctx, info)
	default:
		return false
	}
}

//func (s *meetingServer) refreshRepeatMeeting(ctx context.Context, info *model.MeetingInfo) map[string]any {
//	updateData := map[string]any{}
//	// get current timestamp
//	nowTimestamp := timeutil.GetCurrentTimestampBySecond()
//	if info.EndDate > 0 && info.EndDate < nowTimestamp {
//		updateData["status"] = constant.Completed
//		return updateData
//	}
//	// get days between start time and current time
//	// info.ScheduledTime + info.MeetingDuration
//	if info.RepeatTimes > 0 {
//		days, err := timeutil.DaysBetweenTimestamps(info.TimeZone, info.ScheduledTime+info.MeetingDuration)
//		if err != nil {
//			return updateData
//		}
//		if int32(days) > info.RepeatTimes {
//			updateData["status"] = constant.Completed
//		}
//		return updateData
//	}
//
//	if !s.IsTodayNeedMeeting(ctx, info) {
//		updateData["status"] = constant.Scheduled
//		return updateData
//	}
//
//	startTime := s.GetDayTimestamp(info.ScheduledTime)
//	now := s.GetDayTimestamp(nowTimestamp)
//	if startTime+info.MeetingDuration < now {
//		updateData["status"] = constant.Scheduled
//	} else if startTime+info.MeetingDuration < now && startTime > now {
//		if info.Status == constant.InProgress {
//			return updateData
//		}
//		updateData["status"] = constant.InProgress
//	}
//
//	return updateData
//}

func (s *meetingServer) refreshRepeatMeeting(ctx context.Context, info *model.MeetingInfo) map[string]any {
	updateData := map[string]any{}
	// get current timestamp
	nowTimestamp := timeutil.GetCurrentTimestampBySecond()
	if info.EndDate > 0 && info.EndDate < nowTimestamp {
		updateData["status"] = constant.Completed
		return updateData
	}

	if !s.IsTodayNeedMeeting(ctx, info) {
		updateData["status"] = constant.Scheduled
		nextTime := s.nextMeetingTimestamp(ctx, info)
		if nextTime == 0 {
			updateData["status"] = constant.Completed
		}
		return updateData
	}

	startTime := s.GetDayTimestamp(info.ScheduledTime)
	now := s.GetDayTimestamp(nowTimestamp)
	if startTime+info.MeetingDuration < now {
		updateData["status"] = constant.Scheduled
	} else if startTime+info.MeetingDuration < now && startTime > now {
		if info.Status == constant.InProgress {
			return updateData
		}
		updateData["status"] = constant.InProgress
	}

	return updateData
}

func (s *meetingServer) GetDayTimestamp(timestamp int64) int64 {
	return timestamp - timeutil.GetCurDayZeroTimestamp()
}

func (s *meetingServer) isValidDay(day time.Weekday, days []int32) bool {
	for _, d := range days {
		if int32(day) == d {
			return true
		}
	}
	return false
}

func (s *meetingServer) nextMeetingTimestamp(ctx context.Context, info *model.MeetingInfo) int64 {
	loc, err := time.LoadLocation(info.TimeZone)
	if err != nil {
		log.ZError(ctx, "load location failed", err, "timezone invalid", info.TimeZone)
		fmt.Println("Unable to parse timezone:", err)
		return 0
	}

	scheduledTime := time.Unix(info.ScheduledTime, 0).In(loc)
	now := time.Now().In(loc)

	if info.RepeatType == constant.NoneRepeat {
		if now.Before(scheduledTime) {
			return scheduledTime.Unix()
		}
		return 0
	}

	if info.EndDate > 0 && now.After(time.Unix(info.EndDate, 0).In(loc)) {
		return 0
	}

	var nextTime time.Time

	if now.After(scheduledTime) {
		switch info.RepeatType {
		case constant.RepeatDaily:
			for nextTime = scheduledTime; !now.Before(nextTime); nextTime = nextTime.AddDate(0, 0, 1) {
			}
		case constant.RepeatWeekly:
			for nextTime = scheduledTime; !now.Before(nextTime); nextTime = nextTime.AddDate(0, 0, 7) {
			}
		case constant.RepeatWeekDay:
			nextTime = scheduledTime
			for !now.Before(nextTime) || nextTime.Weekday() == time.Saturday || nextTime.Weekday() == time.Sunday {
				nextTime = nextTime.AddDate(0, 0, 1)
			}
		case constant.RepeatMonth:
			for nextTime = scheduledTime; !now.Before(nextTime); nextTime = nextTime.AddDate(0, 1, 0) {
			}
		case constant.RepeatCustom:
			nextTime = scheduledTime
			interval := int(info.Interval)
			for !now.Before(nextTime) {
				switch info.UintType {
				case constant.UnitTypeDay:
					nextTime = nextTime.AddDate(0, 0, interval)
				case constant.UnitTypeWeek:
					nextTime = nextTime.AddDate(0, 0, interval*7)
				case constant.UnitTypeMonth:
					nextTime = nextTime.AddDate(0, interval, 0)
				}
				if len(info.RepeatDayOfWeek) > 0 {
					for !s.isValidDay(nextTime.Weekday(), info.RepeatDayOfWeek) {
						nextTime = nextTime.AddDate(0, 0, 1)
					}
				}
			}
		}
	} else {
		nextTime = scheduledTime
	}

	if info.EndDate > 0 && nextTime.After(time.Unix(info.EndDate, 0).In(loc)) {
		return 0
	}

	return nextTime.Unix()
}

func (s *meetingServer) refreshMeetingStatus(ctx context.Context) {
	meetings, err := s.meetingStorageHandler.FindByStatus(ctx, []string{constant.InProgress, constant.Scheduled}, "")
	if err != nil {
		log.ZError(ctx, "find meetings failed", err)
		return
	}
	for _, one := range meetings {
		updateData := map[string]any{}
		if one.RepeatType == constant.NoneRepeat || one.RepeatType == "" {
			updateData = s.refreshNonRepeatMeeting(ctx, one)
		} else {
			updateData = s.refreshRepeatMeeting(ctx, one)
		}
		if len(updateData) > 0 {
			if updateData["status"] == constant.Completed {
				_ = s.handleCompleteMeeting(ctx, one.MeetingID)
			}
			if err := s.meetingStorageHandler.Update(ctx, one.MeetingID, updateData); err != nil {
				log.ZError(ctx, "update meeting status failed", err)
			}
		}
	}
}

func (s *meetingServer) checkCanStartMeeting(info *model.MeetingInfo) bool {
	now := timeutil.GetCurrentTimestampBySecond()
	if info.RepeatType == "" && info.ScheduledTime+info.MeetingDuration < now {
		return false
	}
	return true
}

func (s *meetingServer) handleCompleteMeeting(ctx context.Context, meetingID string) error {
	if err := s.meetingRtc.CloseRoom(ctx, meetingID); err != nil {
		log.ZError(ctx, "handle complete meeting close room failed", err, "meetingID", meetingID)
		return err
	}
	return nil
}

func (s *meetingServer) toggleEnableMeeting(ctx context.Context, meetingID string) (*pbmeeting.ToggleRecordMeetingResp, error) {
	resp := &pbmeeting.ToggleRecordMeetingResp{}
	egressID, downloadUrl, err := s.meetingRtc.StartUpload(ctx, meetingID)
	if err != nil {
		return resp, servererrs.ErrEnableRecordMeetingFailed.WrapMsg("toggle enable record meeting failed.")
	}
	updateData := map[string]any{
		"egress_id":    egressID,
		"download_url": downloadUrl,
	}
	if err := s.meetingStorageHandler.Update(ctx, meetingID, updateData); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *meetingServer) toggleDisableMeeting(ctx context.Context, egressID string) (*pbmeeting.ToggleRecordMeetingResp, error) {
	resp := &pbmeeting.ToggleRecordMeetingResp{}

	if egressID == "" {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("egress_id is null, no need to disable")
	}

	if err := s.meetingRtc.StopUpload(ctx, egressID); err != nil {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("toggle disable record meeting failed.")
	}
	return resp, nil
}
