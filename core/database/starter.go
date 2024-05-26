package database

import (
	"github.com/gophab/gophrame/core/database/config"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing Database: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		global.DB = InitDB()
		inject.InjectValue("database", global.DB)
	}
}
