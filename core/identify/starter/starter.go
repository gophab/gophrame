package starter

import (
	"sync"

	_ "github.com/gophab/gophrame/core/identify"

	"github.com/gophab/gophrame/core/identify/aliyun"
	"github.com/gophab/gophrame/core/identify/tcloud"

	"github.com/gophab/gophrame/core/identify/config"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Starting Identify: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			aliyun.Start()
			tcloud.Start()
		})
	}
}
