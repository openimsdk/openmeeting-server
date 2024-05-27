package main

import (
	_ "net/http/pprof"

	"github.com/openimsdk/openmeeting-server/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewAdminApiCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}

}
