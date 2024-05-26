package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

var ConfigYml config.IYmlConfig = config.ConfigYml

func Init() {
	logger.Info("Initializing Framework Config...")
	config.InitConfig()
}
