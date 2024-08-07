package aliyun

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/oss/aliyun/config"
)

func Start() {
	if config.Setting.Enabled {
		logger.Info("Start Aliyun OSS...")
		if oss, err := CreateAliyunOSS(); err == nil && oss != nil {
			inject.InjectValue("oss", oss)
		}
	}
}
