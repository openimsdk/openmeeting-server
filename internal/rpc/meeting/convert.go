package meeting

import (
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/openimsdk/openmeeting-server/pkg/common/constant"
	"github.com/openimsdk/openmeeting-server/pkg/common/convert"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	pbuser "github.com/openimsdk/protocol/openmeeting/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/timeutil"
	"reflect"
	"sort"
	"strings"
)

func (s *meetingServer) getHostUserID(metadata *pbmeeting.MeetingMetadata) string {
	return metadata.Detail.Info.CreatorDefinedMeeting.HostUserID
}

func (s *meetingServer) getCreatorUserID(metadata *pbmeeting.MeetingMetadata) string {
	return metadata.Detail.Info.SystemGenerated.CreatorUserID
}

func (s *meetingServer) generateMeetingDBData4Booking(ctx context.Context, req *pbmeeting.BookMeetingReq) (*model.MeetingInfo, error) {
	meetingID, err := s.meetingStorageHandler.GenerateMeetingID(ctx)
	if err != nil {
		return nil, errs.WrapMsg(err, "generate meeting id failed")
	}

	dbInfo := &model.MeetingInfo{
		MeetingID:       meetingID,
		Title:           req.CreatorDefinedMeetingInfo.Title,
		StartTime:       req.CreatorDefinedMeetingInfo.ScheduledTime,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		TimeZone:        req.CreatorDefinedMeetingInfo.TimeZone,
		Status:          constant.Scheduled,
		CreatorUserID:   req.CreatorUserID,
	}
	marshal := jsonpb.Marshaler{}
	setting, err := marshal.MarshalToString(req.Setting)
	if err != nil {
		return nil, errs.WrapMsg(err, "marshal send data failed")
	}
	dbInfo.Setting = setting
	if req.RepeatInfo != nil {
		dbInfo.EndDate = req.RepeatInfo.EndDate
		dbInfo.RepeatTimes = req.RepeatInfo.RepeatTimes
		dbInfo.RepeatType = req.RepeatInfo.RepeatType
		if req.RepeatInfo.RepeatType == constant.RepeatCustom {
			dbInfo.UintType = req.RepeatInfo.UintType
			dbInfo.Interval = req.RepeatInfo.Interval
			if req.RepeatInfo.RepeatDaysOfWeek != nil {
				dbInfo.RepeatDayOfWeek = *s.getDBRepeatDayOfWeek(&req.RepeatInfo.RepeatDaysOfWeek)
			}
		}
	}
	return dbInfo, nil
}

func (s *meetingServer) getDBRepeatDayOfWeek(weeks *[]pbmeeting.DayOfWeek) *[]int32 {
	var dayOfWeek []int32
	for _, one := range *weeks {
		dayOfWeek = append(dayOfWeek, int32(one))
	}
	sort.Sort(convert.ByInt32(dayOfWeek))
	return &dayOfWeek
}

func (s *meetingServer) getClientRepeatDayOfWeek(dayOfWeek *[]int32) *[]pbmeeting.DayOfWeek {
	days := &[]pbmeeting.DayOfWeek{}
	for _, one := range *dayOfWeek {
		*days = append(*days, pbmeeting.DayOfWeek(one))
	}
	return days
}

func (s *meetingServer) generateMeetingDBData4Create(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*model.MeetingInfo, error) {
	meetingID, err := s.meetingStorageHandler.GenerateMeetingID(ctx)
	if err != nil {
		return nil, errs.WrapMsg(err, "generate meeting id failed")
	}

	marshal := jsonpb.Marshaler{}
	settingString, err := marshal.MarshalToString(req.Setting)
	if err != nil {
		return nil, errs.WrapMsg(err, "marshal send data failed")
	}

	return &model.MeetingInfo{
		MeetingID:       meetingID,
		Title:           req.CreatorDefinedMeetingInfo.Title,
		StartTime:       timeutil.GetCurrentTimestampBySecond(),
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		Status:          constant.InProgress,
		CreatorUserID:   req.CreatorUserID,
		Setting:         settingString,
	}, nil
}

func (s *meetingServer) generateParticipantMetaData(userInfo *pbuser.GetUserInfoResp) *pbmeeting.ParticipantMetaData {
	return &pbmeeting.ParticipantMetaData{
		UserInfo: &pbmeeting.UserInfo{
			UserID:   userInfo.UserID,
			Nickname: userInfo.Nickname,
			Account:  userInfo.Account,
		},
	}
}

func (s *meetingServer) generateDefaultPersonalData(userID string) *pbmeeting.PersonalData {
	return &pbmeeting.PersonalData{
		UserID: userID,
		PersonalSetting: &pbmeeting.PersonalMeetingSetting{
			CameraOnEntry:     false,
			MicrophoneOnEntry: false,
		},
		LimitSetting: &pbmeeting.PersonalMeetingSetting{
			CameraOnEntry:     true,
			MicrophoneOnEntry: true,
		},
	}
}

func (s *meetingServer) generateMeetingMetaData(ctx context.Context, info *model.MeetingInfo) (*pbmeeting.MeetingMetadata, error) {
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: info.CreatorUserID})
	if err != nil {
		return nil, errs.WrapMsg(err, "get user info failed")
	}

	metaData := &pbmeeting.MeetingMetadata{}
	metaData.PersonalData = []*pbmeeting.PersonalData{s.generateDefaultPersonalData(info.CreatorUserID)}
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID:   info.CreatorUserID,
		Status:          info.Status,
		StartTime:       info.StartTime,
		MeetingID:       info.MeetingID,
		CreatorNickname: userInfo.Nickname,
	}
	creatorInfo := &pbmeeting.CreatorDefinedMeetingInfo{
		Title:           info.Title,
		ScheduledTime:   info.ScheduledTime,
		MeetingDuration: info.MeetingDuration,
		Password:        info.Password,
		TimeZone:        info.TimeZone,
		HostUserID:      info.CreatorUserID,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	setting := &pbmeeting.MeetingSetting{}
	if info.Setting != "" {
		unMarshal := jsonpb.Unmarshaler{}
		if err := unMarshal.Unmarshal(strings.NewReader(info.Setting), setting); err != nil {
			return nil, errs.WrapMsg(err, "unMarshal db data failed")
		}
	}

	repeatInfo := &pbmeeting.MeetingRepeatInfo{
		RepeatType:       info.RepeatType,
		EndDate:          info.EndDate,
		RepeatTimes:      info.RepeatTimes,
		UintType:         info.UintType,
		Interval:         info.Interval,
		RepeatDaysOfWeek: *s.getClientRepeatDayOfWeek(&info.RepeatDayOfWeek),
	}

	metaData.Detail = &pbmeeting.MeetingInfoSetting{
		Setting:    setting,
		Info:       meetingInfo,
		RepeatInfo: repeatInfo,
	}
	return metaData, nil
}

func (s *meetingServer) getMeetingDetailSetting(ctx context.Context, info *model.MeetingInfo) (*pbmeeting.MeetingInfoSetting, error) {
	// Fill in response data
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: info.CreatorUserID})
	if err != nil {
		log.ZError(ctx, "get user info failed", err, "userID", info.CreatorUserID)
		return nil, errs.WrapMsg(err, "get user info failed")
	}

	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID:   info.CreatorUserID,
		Status:          info.Status,
		StartTime:       info.StartTime,
		MeetingID:       info.MeetingID,
		CreatorNickname: userInfo.Nickname,
	}
	creatorInfo := &pbmeeting.CreatorDefinedMeetingInfo{
		Title:           info.Title,
		ScheduledTime:   info.ScheduledTime,
		MeetingDuration: info.MeetingDuration,
		Password:        info.Password,
		TimeZone:        info.TimeZone,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	repeatInfo := &pbmeeting.MeetingRepeatInfo{
		EndDate:          info.EndDate,
		RepeatTimes:      info.RepeatTimes,
		RepeatType:       info.RepeatType,
		UintType:         info.UintType,
		Interval:         info.Interval,
		RepeatDaysOfWeek: *s.getClientRepeatDayOfWeek(&info.RepeatDayOfWeek),
	}

	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info:       meetingInfo,
		RepeatInfo: repeatInfo,
	}
	if info.Setting != "" {
		setting := &pbmeeting.MeetingSetting{}
		unMarshal := jsonpb.Unmarshaler{}
		if err := unMarshal.Unmarshal(strings.NewReader(info.Setting), setting); err != nil {
			return nil, errs.WrapMsg(err, "unMarshal db data failed")
		}
		meetingInfoSetting.Setting = setting
	}
	meetingInfoSetting.Info.SystemGenerated.CreatorNickname = userInfo.Nickname
	meetingInfoSetting.Info.CreatorDefinedMeeting.MeetingDuration = info.MeetingDuration
	meetingInfoSetting.Info.CreatorDefinedMeeting.HostUserID = info.CreatorUserID

	// first priority
	metaData, err := s.meetingRtc.GetRoomData(ctx, info.MeetingID)
	if err == nil {
		//meetingInfoSetting.Setting = metaData.Detail.Setting
		//meetingInfoSetting.Info.SystemGenerated.CreatorNickname = metaData.Detail.Info.SystemGenerated.CreatorNickname
		//meetingInfoSetting.Info.CreatorDefinedMeeting.MeetingDuration = metaData.Detail.Info.CreatorDefinedMeeting.MeetingDuration
		meetingInfoSetting.Info.CreatorDefinedMeeting.HostUserID = metaData.Detail.Info.CreatorDefinedMeeting.HostUserID
		meetingInfoSetting.Info.CreatorDefinedMeeting.CoHostUSerID = metaData.Detail.Info.CreatorDefinedMeeting.CoHostUSerID
	}

	return meetingInfoSetting, nil
}

func (s *meetingServer) generateMeetingMetaData4Create(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq, info *model.MeetingInfo) (*pbmeeting.MeetingMetadata, error) {
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: info.CreatorUserID})
	if err != nil {
		return nil, errs.WrapMsg(err, "get user info failed")
	}

	metaData := &pbmeeting.MeetingMetadata{}
	metaData.PersonalData = []*pbmeeting.PersonalData{s.generateDefaultPersonalData(req.CreatorUserID)}
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID:   info.CreatorUserID,
		Status:          info.Status,
		StartTime:       info.StartTime,
		MeetingID:       info.MeetingID,
		CreatorNickname: userInfo.Nickname,
	}
	creatorInfo := &pbmeeting.CreatorDefinedMeetingInfo{
		Title:           req.CreatorDefinedMeetingInfo.Title,
		ScheduledTime:   req.CreatorDefinedMeetingInfo.ScheduledTime,
		MeetingDuration: req.CreatorDefinedMeetingInfo.MeetingDuration,
		Password:        req.CreatorDefinedMeetingInfo.Password,
		HostUserID:      req.CreatorUserID,
	}
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: creatorInfo,
	}
	metaData.Detail = &pbmeeting.MeetingInfoSetting{
		Setting: req.Setting,
		Info:    meetingInfo,
	}
	return metaData, nil
}

// generateMeetingInfoSetting generates MeetingInfoSetting from the given request and meeting ID.
func (s *meetingServer) generateClientRespMeetingSetting(
	meetingSetting *pbmeeting.MeetingSetting,
	defineMeetingInfo *pbmeeting.CreatorDefinedMeetingInfo, meeting *model.MeetingInfo) *pbmeeting.MeetingInfoSetting {
	// Fill in response data
	systemInfo := &pbmeeting.SystemGeneratedMeetingInfo{
		CreatorUserID: meeting.CreatorUserID,
		Status:        meeting.Status,
		StartTime:     meeting.ScheduledTime, // Scheduled start time as the actual start time
		MeetingID:     meeting.MeetingID,
	}
	// Combine system-generated and creator-defined info
	meetingInfo := &pbmeeting.MeetingInfo{
		SystemGenerated:       systemInfo,
		CreatorDefinedMeeting: defineMeetingInfo,
	}
	// Create MeetingInfoSetting
	meetingInfoSetting := &pbmeeting.MeetingInfoSetting{
		Info:    meetingInfo,
		Setting: meetingSetting,
	}
	return meetingInfoSetting
}

func (s *meetingServer) getDBUpdateData(ctx context.Context, info *model.MeetingInfo, req *pbmeeting.UpdateMeetingRequest) *map[string]any {
	updateData := map[string]any{}

	// Update main fields
	s.updateMainFields(ctx, &updateData, req)
	log.ZDebug(ctx, "after update main fields", "update data:", updateData)

	// Update repeat info
	s.updateRepeatInfo(&updateData, req)
	log.ZDebug(ctx, "after update repeat info fields", "update data:", updateData)

	// Update settings
	s.updateSettings(info, &updateData, req)
	log.ZDebug(ctx, "after update setting fields", "update data:", updateData)

	return &updateData
}

func (s *meetingServer) updateMainFields(ctx context.Context, updateData *map[string]any, req *pbmeeting.UpdateMeetingRequest) {
	updateField := func(key string, ptr interface{}) {
		val := reflect.ValueOf(ptr)
		if !(val.Kind() == reflect.Ptr && val.IsNil()) {
			v := val.Elem().FieldByName("Value")
			switch v.Kind() {
			case reflect.String:
				(*updateData)[key] = v.String()
			case reflect.Int:
				(*updateData)[key] = v.Int()
			case reflect.Int64:
				(*updateData)[key] = v.Int()
			case reflect.Int32:
				(*updateData)[key] = v.Int()
			case reflect.Bool:
				(*updateData)[key] = v.Bool()

			// Add other types as needed
			default:
				log.ZError(ctx, "Unhandled type", errs.ErrArgs, "type:", v.Kind())
			}
		}
	}

	updateField("title", req.Title)
	updateField("scheduled_time", req.ScheduledTime)
	updateField("meeting_duration", req.MeetingDuration)
	updateField("password", req.Password)
	updateField("time_zone", req.TimeZone)
}

func (s *meetingServer) updateRepeatInfo(updateData *map[string]any, req *pbmeeting.UpdateMeetingRequest) {
	if req.RepeatInfo != nil {
		(*updateData)["end_date"] = req.RepeatInfo.EndDate
		(*updateData)["repeat_times"] = req.RepeatInfo.RepeatTimes
		(*updateData)["repeat_type"] = req.RepeatInfo.RepeatType
		if req.RepeatInfo.RepeatType == constant.RepeatCustom {
			(*updateData)["uint_type"] = req.RepeatInfo.UintType
			(*updateData)["interval"] = req.RepeatInfo.Interval
			(*updateData)["repeat_day_of_week"] = *s.getDBRepeatDayOfWeek(&req.RepeatInfo.RepeatDaysOfWeek)
		} else {
			// reset setting
			(*updateData)["uint_type"] = ""
			(*updateData)["interval"] = 0
			(*updateData)["repeat_day_of_week"] = nil
		}
	}
}

func (s *meetingServer) updateSettings(info *model.MeetingInfo, updateData *map[string]any, req *pbmeeting.UpdateMeetingRequest) {
	setting := &pbmeeting.MeetingSetting{}
	if info.Setting != "" {
		unMarshal := jsonpb.Unmarshaler{}
		if err := unMarshal.Unmarshal(strings.NewReader(info.Setting), setting); err != nil {
			return
		}
	}

	updateSetting := false
	updateSettingField := func(field interface{}, ptr interface{}) {
		val := reflect.ValueOf(ptr)
		if !(val.Kind() == reflect.Ptr && val.IsNil()) {
			updateSetting = true
			v := val.Elem().FieldByName("Value")
			reflect.ValueOf(field).Elem().Set(v)
		}
	}

	updateSettingField(&setting.CanParticipantsEnableCamera, req.CanParticipantsEnableCamera)
	updateSettingField(&setting.CanParticipantsUnmuteMicrophone, req.CanParticipantsUnmuteMicrophone)
	updateSettingField(&setting.CanParticipantsShareScreen, req.CanParticipantsShareScreen)
	updateSettingField(&setting.DisableCameraOnJoin, req.DisableCameraOnJoin)
	updateSettingField(&setting.DisableMicrophoneOnJoin, req.DisableMicrophoneOnJoin)
	updateSettingField(&setting.CanParticipantJoinMeetingEarly, req.CanParticipantJoinMeetingEarly)
	updateSettingField(&setting.AudioEncouragement, req.AudioEncouragement)
	updateSettingField(&setting.LockMeeting, req.LockMeeting)
	updateSettingField(&setting.VideoMirroring, req.VideoMirroring)

	if updateSetting {
		marshal := jsonpb.Marshaler{}
		updateString, err := marshal.MarshalToString(setting)
		if err != nil {
			return
		}
		(*updateData)["setting"] = updateString
	}
}

func (s *meetingServer) mergeAndUnique(array1, array2 []string) []string {
	exists := make(map[string]bool)
	var result []string

	for _, v := range array1 {
		if !exists[v] {
			exists[v] = true
			result = append(result, v)
		}
	}
	for _, v := range array2 {
		if !exists[v] {
			exists[v] = true
			result = append(result, v)
		}
	}
	return result
}
