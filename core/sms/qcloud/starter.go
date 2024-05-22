package qcloud

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/sms/qcloud/config"
)

func Start() {
	if config.Setting.Enabled {
		if sender, err := CreateQcloudSmsSender(); err == nil && sender != nil {
			inject.InjectValue("smsSender", sender)
		}
	}
}
