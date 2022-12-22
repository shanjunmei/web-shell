package main

import (
	"github.com/jiangklijna/web-shell/cmd"
	"github.com/jiangklijna/web-shell/core"
)

func main() {
	core.StartService(service)
}

func service() {
	parms := new(cmd.Parameter)
	parms.Init()
	parms.Run()
}
