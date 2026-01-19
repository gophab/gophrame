package database

import (
	"github.com/gophab/gophrame/core/database"
	"github.com/gophab/gophrame/core/database/config"

	_ "github.com/gophab/gophrame/core/database/mysql"
	_ "github.com/gophab/gophrame/core/database/postgres"

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
		var err error
		if global.DB, err = database.InitDB(); err == nil {
			inject.InjectValue("database", global.DB)
			logger.Info("Database initialized.")
		} else {
			logger.Error("Initializing Database error: ", err.Error())
		}
	}
}
