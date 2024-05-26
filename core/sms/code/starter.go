package code

import (
	"sync"

	"github.com/gophab/gophrame/core/code"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/core/sms/code/config"
	SmsConfig "github.com/gophab/gophrame/core/sms/config"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	once.Do(func() {
		initSmsCodeSender()
		initSmsCodeStore()
		initSmsCodeValidator()
		initSmsCodeController()
	})
}

func initSmsCodeSender() (sender code.CodeSender, err error) {
	if SmsConfig.Setting.Enabled && config.Setting.Enabled {
		sender := &SmsCodeSender{}
		inject.InjectValue("smsCodeSender", sender)
	}
	return sender, err
}

func initSmsCodeStore() (store code.CodeStore, err error) {
	if SmsConfig.Setting.Enabled && config.Setting.Enabled {
		if config.Setting.Redis != nil && config.Setting.Redis.Enabled {
			store, err = code.CreateRedisCodeStore(config.Setting)
		} else if config.Setting.Cache != nil && config.Setting.Cache.Enabled {
			store, err = code.CreateCacheCodeStore(config.Setting)
		} else {
			store, err = code.CreateMemoryCodeStore(config.Setting)
		}
		if store != nil {
			inject.InjectValue("smsCodeStore", store)
		}
	}
	return store, err
}

func initSmsCodeValidator() {
	if SmsConfig.Setting.Enabled && config.Setting.Enabled {
		inject.InjectValue("smsCodeValidator", &SmsCodeValidator{})
	}
}

func initSmsCodeController() {
	if SmsConfig.Setting.Enabled && config.Setting.Enabled {
		smsCodeController := &SmsCodeController{}
		inject.InjectValue("smsCodeController", smsCodeController)

		// 注册Controller
		controller.AddController(smsCodeController)
	}
}
