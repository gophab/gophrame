package starter

import (
	"sync"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/sms/qcloud"
	"github.com/wjshen/gophrame/core/sms/qcloud/config"

	_ "github.com/wjshen/gophrame/config"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		if config.Setting.Enabled {
			if sender, err := qcloud.CreateQcloudSmsSender(); err == nil && sender != nil {
				inject.InjectValue("smsSender", sender)
			}
		}
	})
}
