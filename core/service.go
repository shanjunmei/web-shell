package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

var serviceLogger service.Logger

type app struct {
	work func()
}

func (p *app) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	log.Println("service start ", time.Now())
	go p.work()
	return nil
}
func (p *app) run() {
	log.Println("service run ", time.Now())

	p.work()
}
func (p *app) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	os.Exit(1)
	return nil
}

func PrintHelp() {
	fmt.Println(`Usage:
	web-shell install/uninstall/stop/start
  
  Example:
	web-shell install
	web-shell start
  `)
}

func StartService(work func()) {
	args := os.Args
	if len(args) > 1 {
		controlAction := args[1]
		file, err := os.Executable()
		if err != nil {
			panic(err)
		}
		wd := filepath.Dir(file)
		log.Println("WorkingDirectory :", wd)
		ServiceArg := make([]string, 0)
		ServiceArg[0] = "-s"
		svcConfig := &service.Config{
			Name:             "web-shell",
			DisplayName:      "Web Shell Service",
			Description:      "Web Shell Service",
			Arguments:        ServiceArg,
			WorkingDirectory: wd,
		}

		prg := &app{work: work}
		svc, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		serviceLogger, err = svc.Logger(nil)

		var isControlFlag = false
		for _, ca := range service.ControlAction {
			if ca == controlAction {
				isControlFlag = true
			}
		}
		if isControlFlag {
			err = service.Control(svc, controlAction)
			if err != nil {
				log.Println("service control fail :", err)
			}
		} else {
			PrintHelp()
			work()

		}
	} else {
		work()
	}

}
