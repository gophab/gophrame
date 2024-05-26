package redis

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/redis/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing Redis: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		initRedisClientPool(config.Setting.Database)
	}
}
