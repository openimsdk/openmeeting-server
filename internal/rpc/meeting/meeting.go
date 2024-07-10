package meeting

import (
	"context"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common"
	"github.com/openimsdk/openmeeting-server/pkg/common/constant"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	sysConstant "github.com/openimsdk/protocol/constant"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	pbuser "github.com/openimsdk/protocol/openmeeting/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

// BookMeeting Implement the MeetingServiceServer interface
func (s *meetingServer) BookMeeting(ctx context.Context, req *pbmeeting.BookMeetingReq) (*pbmeeting.BookMeetingResp, error) {
	resp := &pbmeeting.BookMeetingResp{}
	meetingDBInfo, err := s.generateMeetingDBData4Booking(ctx, req)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting data failed")
	}
	metaData, err := s.generateMeetingMetaData(ctx, meetingDBInfo)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting meta data failed")
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}
	resp.Detail = metaData.Detail
	return resp, nil
}

func (s *meetingServer) CreateImmediateMeeting(ctx context.Context, req *pbmeeting.CreateImmediateMeetingReq) (*pbmeeting.CreateImmediateMeetingResp, error) {
	resp := &pbmeeting.CreateImmediateMeetingResp{}

	inMeeting, err := s.checkUserInMeeting(ctx, req.CreatorUserID)
	if err != nil {
		return resp, errs.WrapMsg(err, "create meeting failed")
	}
	if inMeeting {
		return resp, servererrs.ErrMeetingUserLimit.WrapMsg("user already in meeting")
	}

	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.CreatorUserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	meetingDBInfo, err := s.generateMeetingDBData4Create(ctx, req)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting data failed")
	}

	metaData, err := s.generateMeetingMetaData4Create(ctx, req, meetingDBInfo)
	if err != nil {
		return resp, errs.WrapMsg(err, "generate meeting meta data failed")
	}
	participantMetaData := s.generateParticipantMetaData(userInfo)

	_, token, liveUrl, err := s.meetingRtc.CreateRoom(ctx, meetingDBInfo.MeetingID, req.CreatorUserID, metaData, participantMetaData, s.userRpc)
	if err != nil {
		return resp, err
	}

	err = s.meetingStorageHandler.Create(ctx, []*model.MeetingInfo{meetingDBInfo})
	if err != nil {
		return resp, err
	}

	// create meeting meta data
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, err
	}

	resp.Detail = metaData.Detail
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) JoinMeeting(ctx context.Context, req *pbmeeting.JoinMeetingReq) (*pbmeeting.JoinMeetingResp, error) {
	resp := &pbmeeting.JoinMeetingResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.UserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	dbInfo, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get meeting data failed")
	}

	inMeeting, err := s.checkUserInMeeting(ctx, req.UserID)
	if err != nil {
		return resp, errs.WrapMsg(err, "join meeting failed")
	}
	if inMeeting {
		return resp, servererrs.ErrMeetingUserLimit.WrapMsg("user already in meeting")
	}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		if !s.checkCanStartMeeting(dbInfo) {
			return resp, errs.WrapMsg(err, "get room data failed")
		}
		// for those need repeat booking meeting, create new rooms
		metaData, err = s.generateMeetingMetaData(ctx, dbInfo)
		if err != nil {
			return resp, errs.WrapMsg(err, "generate meeting meta data failed")
		}
		participantMetaData := s.generateParticipantMetaData(userInfo)
		if ps, err := s.meetingRtc.ListParticipants(ctx, req.MeetingID); err != nil {
			for _, p := range ps {
				if p.Identity == req.UserID {
					participantMetaData = nil
					break
				}
			}
		}
		_, _, _, err = s.meetingRtc.CreateRoom(ctx, dbInfo.MeetingID, userInfo.UserID, metaData, participantMetaData, s.userRpc)
		if err != nil {
			return resp, err
		}
	}

	if req.UserID != s.getHostUserID(metaData) && req.Password != metaData.Detail.Info.CreatorDefinedMeeting.Password {
		return resp, servererrs.ErrMeetingPasswordNotMatch.WrapMsg("meeting password not match, please check and try again!")
	}

	metaData.Detail.Info.SystemGenerated.MeetingID = req.MeetingID
	participantMetaData := s.generateParticipantMetaData(userInfo)

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.UserID, participantMetaData)
	if err != nil {
		return resp, errs.WrapMsg(err, "get join token failed")
	}

	// update meta data to liveKit
	found := false
	for _, personalData := range metaData.PersonalData {
		if personalData.UserID == req.UserID {
			found = true
			break
		}
	}
	if !found {
		personalData := s.generateDefaultPersonalData(req.UserID)
		metaData.PersonalData = append(metaData.PersonalData, personalData)
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, errs.WrapMsg(err, "update meta data failed")
	}
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) GetMeetingToken(ctx context.Context, req *pbmeeting.GetMeetingTokenReq) (*pbmeeting.GetMeetingTokenResp, error) {
	resp := &pbmeeting.GetMeetingTokenResp{}
	userInfo, err := s.userRpc.Client.GetUserInfo(ctx, &pbuser.GetUserInfoReq{UserID: req.UserID})
	if err != nil {
		return resp, errs.WrapMsg(err, "get user info failed")
	}

	participantMetaData := s.generateParticipantMetaData(userInfo)

	token, liveUrl, err := s.meetingRtc.GetJoinToken(ctx, req.MeetingID, req.UserID, participantMetaData)
	if err != nil {
		return resp, err
	}

	resp.MeetingID = req.MeetingID
	resp.LiveKit = &pbmeeting.LiveKit{
		Token: token,
		Url:   liveUrl,
	}
	return resp, nil
}

func (s *meetingServer) LeaveMeeting(ctx context.Context, req *pbmeeting.LeaveMeetingReq) (*pbmeeting.LeaveMeetingResp, error) {
	resp := &pbmeeting.LeaveMeetingResp{}

	if err := s.meetingRtc.RemoveParticipant(ctx, req.MeetingID, req.UserID); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *meetingServer) EndMeeting(ctx context.Context, req *pbmeeting.EndMeetingReq) (*pbmeeting.EndMeetingResp, error) {
	resp := &pbmeeting.EndMeetingResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}
	if !s.checkAuthPermission(metaData.Detail.Info.CreatorDefinedMeeting.HostUserID, req.UserID) {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("user did not have permission to end somebody's meeting")
	}
	dbInfo, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get meeting data failed")
	}

	// change status to completed
	status := constant.Completed
	// if we have next meeting schedule
	if s.nextMeetingTimestamp(ctx, dbInfo) > 0 {
		status = constant.Scheduled
	}

	updateData := map[string]any{
		"status": status,
	}
	if err := s.meetingRtc.CloseRoom(ctx, req.MeetingID); err != nil {
		return resp, err
	}
	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, updateData); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *meetingServer) GetMeetings(ctx context.Context, req *pbmeeting.GetMeetingsReq) (*pbmeeting.GetMeetingsResp, error) {
	s.refreshMeetingStatus(ctx)
	resp := &pbmeeting.GetMeetingsResp{}
	meetings, err := s.meetingStorageHandler.FindByStatus(ctx, req.Status)
	if err != nil {
		return resp, err
	}

	// Create response
	var meetingDetails []*pbmeeting.MeetingInfoSetting
	for _, meeting := range meetings {
		detailSetting, err := s.getMeetingDetailSetting(ctx, meeting)
		if err != nil {
			return resp, err
		}
		meetingDetails = append(meetingDetails, detailSetting)
	}
	resp.MeetingDetails = meetingDetails
	return resp, nil
}

func (s *meetingServer) GetMeeting(ctx context.Context, req *pbmeeting.GetMeetingReq) (*pbmeeting.GetMeetingResp, error) {
	resp := &pbmeeting.GetMeetingResp{}
	meetingDBInfo, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	detailSetting, err := s.getMeetingDetailSetting(ctx, meetingDBInfo)
	if err != nil {
		return resp, err
	}
	resp.MeetingDetail = detailSetting
	return resp, nil
}

func (s *meetingServer) UpdateMeeting(ctx context.Context, req *pbmeeting.UpdateMeetingRequest) (*pbmeeting.UpdateMeetingResp, error) {
	resp := &pbmeeting.UpdateMeetingResp{}

	info, err := s.meetingStorageHandler.TakeWithError(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	if info.Status == constant.Completed {
		log.CInfo(ctx, "meeting is already completed, can not update anymore", "meetingID:", req.MeetingID)
		return resp, servererrs.ErrMeetingAlreadyCompleted.WrapMsg("meeting is already completed, can not update anymore")
	}

	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		//return resp, err
		log.CInfo(ctx, "not found room info in livekit", "meetingID:", req.MeetingID)
	}

	updateData := s.getDBUpdateData(ctx, info, req)

	if len(*updateData) == 0 {
		log.ZDebug(ctx, "no need to update meeting", "meeting id", req.MeetingID)
		return resp, nil
	}
	if err := s.meetingStorageHandler.Update(ctx, req.MeetingID, *updateData); err != nil {
		return resp, err
	}

	// do not get the metadata, then return successfully
	if metaData == nil {
		return resp, nil
	}

	// get the latest data from database and update metadata
	if err := s.updateMeetingMetaData(ctx, req.MeetingID, metaData); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *meetingServer) updateMeetingMetaData(ctx context.Context, meetingID string, metaData *pbmeeting.MeetingMetadata) error {
	info, err := s.meetingStorageHandler.TakeWithError(ctx, meetingID)
	if err != nil {
		return err
	}

	metaData.Detail, err = s.getMeetingDetailSetting(ctx, info)
	if err != nil {
		return err
	}

	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return err
	}

	return nil
}

func (s *meetingServer) GetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.GetPersonalMeetingSettingsReq) (*pbmeeting.GetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.GetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}
	for _, personalData := range metaData.PersonalData {
		if personalData.GetUserID() == req.UserID {
			resp.Setting = personalData.PersonalSetting
			break
		}
	}
	return resp, nil
}

func (s *meetingServer) SetPersonalMeetingSettings(ctx context.Context, req *pbmeeting.SetPersonalMeetingSettingsReq) (*pbmeeting.SetPersonalMeetingSettingsResp, error) {
	resp := &pbmeeting.SetPersonalMeetingSettingsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	result, err := common.GetKeyFromContext(ctx, sysConstant.OpUserID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get userid from context failed")
	}
	userID := result.(string)
	hostUser := s.getHostUserID(metaData)
	// non host user can not change other people's personal setting
	if userID != hostUser && userID != req.UserID {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("do not have the permission to change other participant's setting")
	}

	if userID == req.UserID {
		// user set personal setting
		if err := s.setSelfPersonalSetting(ctx, metaData, req); err != nil {
			return resp, errs.WrapMsg(err, "set self personal setting failed")
		}
		return resp, nil
	}
	// admin set participant setting
	if err := s.setParticipantPersonalSetting(ctx, metaData, req); err != nil {
		return resp, errs.WrapMsg(err, "set participant personal setting failed")
	}
	return resp, nil
}

func (s *meetingServer) OperateRoomAllStream(ctx context.Context, req *pbmeeting.OperateRoomAllStreamReq) (*pbmeeting.OperateRoomAllStreamResp, error) {
	resp := &pbmeeting.OperateRoomAllStreamResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, err
	}

	hostUser := s.getHostUserID(metaData)
	if !s.checkAuthPermission(hostUser, req.OperatorUserID) {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("do not have the permission")
	}
	if req.MicrophoneOnEntry != nil {
		resp.StreamNotExistUserIDList, resp.FailedUserIDList, err = s.muteAllStream(ctx, req.MeetingID, audio, !req.MicrophoneOnEntry.Value)
		if err != nil {
			return resp, errs.WrapMsg(err, "operate room all microphone stream failed")
		}
	}

	if req.CameraOnEntry != nil {
		resp.StreamNotExistUserIDList, resp.StreamNotExistUserIDList, err = s.muteAllStream(ctx, req.MeetingID, video, !req.CameraOnEntry.Value)
		if err != nil {
			return resp, errs.WrapMsg(err, "operate room all camera stream failed")
		}
	}

	if err := s.broadcastStreamOperateData(ctx, req, resp.StreamNotExistUserIDList, resp.FailedUserIDList); err != nil {
		return resp, errs.WrapMsg(err, "send notification to all participant failed")
	}

	return resp, nil
}

// ModifyMeetingParticipantNickName modify meeting participant nickname
func (s *meetingServer) ModifyMeetingParticipantNickName(ctx context.Context, req *pbmeeting.ModifyMeetingParticipantNickNameReq) (*pbmeeting.ModifyMeetingParticipantNickNameResp, error) {
	resp := &pbmeeting.ModifyMeetingParticipantNickNameResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get room data failed", req.MeetingID)
	}
	// check permission
	if !s.checkAuthPermission(metaData.Detail.Info.CreatorDefinedMeeting.HostUserID, req.UserID) {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("user did not have permission to modify meeting participant's nickname")
	}
	participantMetaData, err := s.meetingRtc.GetParticipantMetaData(ctx, req.MeetingID, req.ParticipantUserID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get participant data failed")
	}
	participantMetaData.UserInfo.Nickname = req.Nickname
	if err = s.meetingRtc.UpdateParticipantData(ctx, participantMetaData, req.MeetingID, req.ParticipantUserID); err != nil {
		return resp, errs.WrapMsg(err, "update participant data failed")
	}
	return resp, nil
}

// RemoveParticipants batch remove participant out of the meeting room
func (s *meetingServer) RemoveParticipants(ctx context.Context, req *pbmeeting.RemoveMeetingParticipantsReq) (*pbmeeting.RemoveMeetingParticipantsResp, error) {
	resp := &pbmeeting.RemoveMeetingParticipantsResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get room data failed", req.MeetingID)
	}
	// check permission only host can remove somebody
	if !s.checkAuthPermission(metaData.Detail.Info.CreatorDefinedMeeting.HostUserID, req.UserID) {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("user did not have permission to remove participant out of the meeting")
	}
	var failedList []string
	var successList []string
	for _, one := range req.ParticipantUserIDs {
		if err = s.meetingRtc.RemoveParticipant(ctx, req.MeetingID, one); err != nil {
			log.ZError(ctx, "remove participant out of the meeting failed", err)
			failedList = append(failedList, one)
		} else {
			successList = append(successList, one)
		}
	}
	resp.FailedUserIDList = failedList
	resp.SuccessUserIDList = successList

	return resp, nil
}

// SetMeetingHostInfo modify host or co-host of the meeting room
func (s *meetingServer) SetMeetingHostInfo(ctx context.Context, req *pbmeeting.SetMeetingHostInfoReq) (*pbmeeting.SetMeetingHostInfoResp, error) {
	resp := &pbmeeting.SetMeetingHostInfoResp{}
	metaData, err := s.meetingRtc.GetRoomData(ctx, req.MeetingID)
	if err != nil {
		return resp, errs.WrapMsg(err, "get room data failed", req.MeetingID)
	}
	// check permission only host can remove somebody
	if !s.checkAuthPermission(metaData.Detail.Info.CreatorDefinedMeeting.HostUserID, req.UserID) {
		return resp, servererrs.ErrMeetingAuthCheck.WrapMsg("user did not have permission to set host info of the meeting")
	}
	if req.HostUserID != nil {
		metaData.Detail.Info.CreatorDefinedMeeting.HostUserID = req.HostUserID.Value
		if err := s.sendMeetingHostData2Client(ctx, req.MeetingID, req.UserID, req.HostUserID.Value, constant.HostTypeHost); err != nil {
			return resp, errs.ErrArgs.WrapMsg("notify host info to participant failed")
		}
	}
	if req.CoHostUserIDs != nil {
		metaData.Detail.Info.CreatorDefinedMeeting.CoHostUSerID = s.mergeAndUnique(
			metaData.Detail.Info.CreatorDefinedMeeting.CoHostUSerID, req.CoHostUserIDs)

		for _, one := range req.CoHostUserIDs {
			if err := s.sendMeetingHostData2Client(ctx, req.MeetingID, req.UserID, one, constant.HostTypeCoHost); err != nil {
				return resp, errs.ErrArgs.WrapMsg("notify host info to participant failed")
			}
		}
	}
	if err := s.meetingRtc.UpdateMetaData(ctx, metaData); err != nil {
		return resp, errs.WrapMsg(err, "update meta data failed")
	}
	return resp, nil
}

func (s *meetingServer) CleanPreviousMeetings(ctx context.Context, req *pbmeeting.CleanPreviousMeetingsReq) (*pbmeeting.CleanPreviousMeetingsResp, error) {
	resp := &pbmeeting.CleanPreviousMeetingsResp{}

	rooms, err := s.meetingRtc.GetAllRooms(ctx)
	if err != nil {
		log.ZError(ctx, "get all rooms error", err, "login and clean previous rooms", req.UserID)
		return resp, errs.WrapMsg(err, "get all rooms error", "login and clean previous rooms", req.UserID)
	}

	for _, room := range rooms {
		ps, err := s.meetingRtc.ListParticipants(ctx, room.Name)
		if err != nil {
			log.ZError(ctx, "list participant error", err, "login and clean previous rooms", room.Name, req.UserID)
		}
		for _, p := range ps {
			fmt.Println(p.Identity, req.UserID)
			if p.Identity != req.UserID {
				continue
			}
			if err := s.notifyKickOffMeetingInfo2Client(ctx, room.Name, req.UserID, constant.KickOffDuplicatedLogin, pbmeeting.KickOffReason_DuplicatedLogin); err != nil {
				log.ZError(ctx, "notify kickoff msg to client error", err, "login and clean previous rooms", room.Name, req.UserID)
			}
			if err := s.meetingRtc.RemoveParticipant(ctx, room.Name, req.UserID); err != nil {
				log.ZError(ctx, "remove participant error", err, "login and clean previous rooms", room.Name, req.UserID)
				continue
			}
		}
	}

	return resp, nil
}
