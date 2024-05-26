package controller

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/router"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Starting Core Controller ...")
	InitRouter(router.Root())
}
