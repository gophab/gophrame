package qcloud

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/oss/qcloud/config"
)

func Start() {
	if config.Setting.Enabled {
		logger.Info("Start Qcloud OSS...")
		if oss, err := CreateQcloudOSS(); err == nil && oss != nil {
			inject.InjectValue("oss", oss)
		}
	}
}
