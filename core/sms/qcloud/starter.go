package qcloud

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/sms/qcloud/config"
)

func Start() {
	if config.Setting.Enabled {
		if sender, err := CreateQcloudSmsSender(); err == nil && sender != nil {
			inject.InjectValue("smsSender", sender)
		}
	}
}
