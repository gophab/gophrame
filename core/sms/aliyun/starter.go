package aliyun

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/sms/aliyun/config"
)

func Start() {
	if config.Setting.Enabled {
		if sender, err := CreateAliyunSmsSender(); err == nil && sender != nil {
			inject.InjectValue("smsSender", sender)
		}
	}
}
