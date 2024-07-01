package oss

import (
	"sync"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/oss/aliyun"
	"github.com/gophab/gophrame/core/oss/config"
	"github.com/gophab/gophrame/core/oss/qcloud"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once          sync.Once
	ossController *OssController
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Enable OSS: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			ossController = &OssController{}
			inject.InjectValue("ossController", ossController)

			aliyun.Start()
			qcloud.Start()

			controller.AddController(ossController)
		})
	}
}
