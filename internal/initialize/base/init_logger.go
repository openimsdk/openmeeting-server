package base

import (
	"github.com/OpenIMSDK/tools/log"
	config "openmeeting-server/dto"
)

func InitLogger() error {
	logConf := config.Config.Log
	if err := log.InitFromConfig("open-meeting.rtc", "rtc", *logConf.RemainLogLevel,
		*logConf.IsStdout, *logConf.IsJson, *logConf.StorageLocation, *logConf.RemainRotationCount,
		*logConf.RotationTime); err != nil {
		panic(err)
	}
	return nil
}
