package starter

import (
	"sync"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"

	"github.com/gophab/gophrame/core/oss"
	"github.com/gophab/gophrame/core/oss/config"

	_ "github.com/gophab/gophrame/core/oss/aliyun"
	_ "github.com/gophab/gophrame/core/oss/qcloud"
)

var (
	once          sync.Once
	ossController *oss.OssController
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Enable OSS: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			ossController = &oss.OssController{}
			inject.InjectValue("ossController", ossController)
			controller.AddController(ossController)
		})
	}
}
