package livekit

import (
	"context"
	"fmt"
	"github.com/livekit/protocol/livekit"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/timeutil"
)

func (x *LiveKit) StartUpload(ctx context.Context, roomID string) (egressID, downloadURL string, err error) {
	egressID, downloadURL, err = "", "", nil
	if !x.uploadConf.Record.Enable {
		log.ZInfo(ctx, "server do not enable record meeting, please check", "roomID", roomID)
		err = servererrs.ErrMeetingRecordSwitchNotOpen.WrapMsg("server do not enable record meeting, please check")
		return
	}
	filename := x.getUploadFilename(roomID)
	downloadURL = x.getDownloadURL(filename)
	request := x.generateRoomCompositeEgressRequest(roomID, filename)
	egressInfo, err := x.egressClient.StartRoomCompositeEgress(ctx, request)
	if err != nil {
		log.ZError(ctx, "", err, "roomID", roomID)
		return
	}
	egressID = egressInfo.EgressId
	return
}

func (u *LiveKit) StopUpload(ctx context.Context, egressID string) error {
	egressInfo, err := u.egressClient.StopEgress(ctx, &livekit.StopEgressRequest{EgressId: egressID})
	if err != nil {
		log.ZError(ctx, "failed to stop recrod meeting", err, "egressID", egressID)
		return err
	}
	log.CInfo(ctx, "stopped record meeting successfully", "roomID", egressInfo.RoomName)
	return nil
}

func (u *LiveKit) getUploadFilename(roomID string) string {
	nowTime := timeutil.GetCurrentTimestampBySecond()
	return "meeting_" + roomID + "_" + fmt.Sprintf("%s", nowTime) + ".mp4"
}

func (u *LiveKit) getDownloadURL(filename string) string {
	return u.uploadConf.S3.Endpoint + "/" + u.uploadConf.S3.Bucket + "/" + filename
}

func (u *LiveKit) generateRoomCompositeEgressRequest(roomID, filename string) *livekit.RoomCompositeEgressRequest {
	encodedFileOutput := &livekit.EncodedFileOutput{
		FileType: livekit.EncodedFileType_MP4,
		Filepath: filename,
		Output: &livekit.EncodedFileOutput_S3{
			S3: &livekit.S3Upload{
				AccessKey: u.uploadConf.S3.AccessKey,
				Secret:    u.uploadConf.S3.Secret,
				Bucket:    u.uploadConf.S3.Bucket,
				Endpoint:  u.uploadConf.S3.Endpoint,
				Region:    u.uploadConf.S3.Region,
			},
		},
	}
	return &livekit.RoomCompositeEgressRequest{
		RoomName:    roomID,
		Layout:      u.uploadConf.Layout,
		FileOutputs: []*livekit.EncodedFileOutput{encodedFileOutput},
		Options: &livekit.RoomCompositeEgressRequest_Advanced{
			Advanced: &livekit.EncodingOptions{
				Framerate: u.uploadConf.FrameRate,
			},
		},
	}
}
