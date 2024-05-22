package database

import (
	"github.com/wjshen/gophrame/core/database/config"
	"github.com/wjshen/gophrame/core/global"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/starter"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Info("Initializing Database")
	if config.Setting.Enabled {
		logger.Info("Database Enabled")
		global.DB = InitDB()
		inject.InjectValue("database", global.DB)
	} else {
		logger.Info("Database Disabled")
	}
}
