package oss

import (
	"sync"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/oss/aliyun"
	"github.com/gophab/gophrame/core/oss/config"
	"github.com/gophab/gophrame/core/oss/qcloud"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Enable OSS: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			aliyun.Start()
			qcloud.Start()
		})
	}
}
