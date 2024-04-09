package response_util

import (
	"openmeeting-server/constant"
	"openmeeting-server/internal/infrastructure/sys_error"
	"openmeeting-server/protocol/pb"
)

func GetSuccessResponseHeader() *pb.RespHeader {
	successCode := constant.SuccessCode
	successMessage := constant.RetCode2MsgMapper[successCode]
	return &pb.RespHeader{
		Retcode: successCode,
		Message: successMessage,
	}
}

func GetErrorResponseHeader(err *sys_error.SysError) *pb.RespHeader {
	return &pb.RespHeader{
		Retcode: err.RetCode,
		Message: err.Msg,
	}
}
