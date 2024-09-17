package cron

import (
	"github.com/gophab/gophrame/core/starter"
	"github.com/robfig/cron/v3"
)

var globalCron = cron.New(cron.WithChain(
	cron.SkipIfStillRunning(cron.DefaultLogger),
))

func init() {
	starter.RegisterStarter(Start)
	starter.RegisterTerminater(Stop)
}

func Start() {
	globalCron.Start()
}

func Stop() {
	globalCron.Stop()
}

func AddFunc(spec string, cmd func()) error {
	_, err := globalCron.AddFunc(spec, cmd)
	return err
}
