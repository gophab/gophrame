package starter

import (
	"sync"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/sms/aliyun"
	"github.com/wjshen/gophrame/core/sms/aliyun/config"

	_ "github.com/wjshen/gophrame/config"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		if config.Setting.Enabled {
			if sender, err := aliyun.CreateAliyunSmsSender(); err == nil && sender != nil {
				inject.InjectValue("smsSender", sender)
			}
		}
	})
}
