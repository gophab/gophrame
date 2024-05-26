package sms

import (
	"sync"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/sms/aliyun"
	"github.com/gophab/gophrame/core/sms/config"
	"github.com/gophab/gophrame/core/sms/qcloud"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Enable SMS: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			aliyun.Start()
			qcloud.Start()
		})
	}
}
