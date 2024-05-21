package starter

import (
	"sync"

	"github.com/wjshen/gophrame/core/code"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/sms/config"

	SmsCode "github.com/wjshen/gophrame/core/sms/code"

	_ "github.com/wjshen/gophrame/config"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		initSmsCodeSender()
		initSmsCodeStore()
		initSmsCodeValidator()
		initSmsCodeController()
	})
}

func initSmsCodeSender() (sender code.CodeSender, err error) {
	if config.Setting.Enabled {
		sender := &SmsCode.SmsCodeSender{}
		inject.InjectValue("smsCodeSender", sender)
	}
	return sender, err
}

func initSmsCodeStore() (store code.CodeStore, err error) {
	if config.Setting.Enabled {
		if config.Setting.Store.Redis != nil && config.Setting.Store.Redis.Enabled {
			store, err = code.CreateRedisCodeStore(config.Setting.Store)
		} else if config.Setting.Store.Cache != nil && config.Setting.Store.Cache.Enabled {
			store, err = code.CreateCacheCodeStore(config.Setting.Store)
		} else {
			store, err = code.CreateMemoryCodeStore(config.Setting.Store)
		}
		if store != nil {
			inject.InjectValue("smsCodeStore", store)
		}
	}
	return store, err
}

func initSmsCodeValidator() {
	if config.Setting.Enabled {
		inject.InjectValue("smsCodeValidator", &SmsCode.SmsCodeValidator{})
	}
}

func initSmsCodeController() {
	if config.Setting.Enabled {
		inject.InjectValue("smsCodeController", &SmsCode.SmsCodeController{})
	}
}
