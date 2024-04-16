package main

import (
	"flag"
	config "openmeeting-server/dto"
	"openmeeting-server/internal/initialize"
	startrpc "openmeeting-server/pkg/common"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "c", "", "config file path")
	flag.Parse()
	if err := config.Parse(confPath); err != nil {
		panic(err)
	}
	if err := startrpc.Start(config.Config.RPC.RTC.Port[0], config.Config.RPC.RTC.Name, initialize.InitServer); err != nil {
		panic(err)
	}
}
