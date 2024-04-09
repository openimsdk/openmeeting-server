package sys_error

import (
	"fmt"
	"openmeeting-server/constant"
)

type SysError struct {
	RetCode int32  `json:"ret_code"`
	Msg     string `json:"msg"`
}

func (s *SysError) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ret_code:%d, msg:%s", s.RetCode, s.Msg)
}

func NewSysError(retcode int32, msg string) *SysError {
	return &SysError{RetCode: retcode, Msg: msg}
}

func NewSysErrorf(retcode int32, format string, a ...interface{}) *SysError {
	return &SysError{RetCode: retcode, Msg: fmt.Sprintf(format, a...)}
}

func HasCodeMess(errorCode int32) bool {
	_, existed := constant.RetCode2MsgMapper[errorCode]
	return existed
}

func GetCodeMess(errorCode int32) string {
	if _, existed := constant.RetCode2MsgMapper[errorCode]; existed {
		return constant.RetCode2MsgMapper[errorCode]
	}
	return ""
}
