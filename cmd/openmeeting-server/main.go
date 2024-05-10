package main

import "openmeeting-server/pkg/common/cmd"
import "github.com/openimsdk/tools/system/program"

func main() {
	if err := cmd.NewMeetingRpcCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
