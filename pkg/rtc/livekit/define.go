package livekit

import (
	lksdk "github.com/livekit/server-sdk-go"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
)

type LiveKit struct {
	roomClient   *lksdk.RoomServiceClient
	egressClient *lksdk.EgressClient
	uploadConf   *config.Upload
	index        uint64
	conf         *config.RTC
}
