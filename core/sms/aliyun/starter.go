package aliyun

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/sms/aliyun/config"
)

func Start() {
	if config.Setting.Enabled {
		if sender, err := CreateAliyunSmsSender(); err == nil && sender != nil {
			inject.InjectValue("smsSender", sender)
		}
	}
}
