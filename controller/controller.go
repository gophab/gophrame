package controller

import (
	"github.com/gophab/gophrame/config"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/webservice/middleware"

	_ "github.com/gophab/gophrame/controller/security"
)

func init() {
	starter.RegisterInitializor(Init)
	starter.RegisterStarter(Start)
}

// Auto Initialize entrypoint
func Init() {
}

// Autostart entrypoint
func Start() {
	logger.Info("Starting Framework Controllers...")
	logger.Info("Allow Cross Domain: ", config.Config.Server.AllowCrossDomain)
	if config.Config.Server.AllowCrossDomain {
		middleware.UseCors()
	}
}
