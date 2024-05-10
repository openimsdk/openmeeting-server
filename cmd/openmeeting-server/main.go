package main

import (
	"github.com/openimsdk/tools/system/program"
	"openmeeting-server/pkg/common/cmd"
)

func main() {
	if err := cmd.NewMeetingRpcCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
