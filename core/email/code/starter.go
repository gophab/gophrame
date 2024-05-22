package code

import (
	"sync"

	"github.com/wjshen/gophrame/core/code"
	"github.com/wjshen/gophrame/core/email/config"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	once.Do(func() {
		initEmailCodeSender()
		initEmailCodeStore()
		initEmailCodeValidator()
		initEmailCodeController()
	})
}

func initEmailCodeSender() (sender code.CodeSender, err error) {
	if config.Setting.Enabled {
		sender = &EmailCodeSender{}
		inject.InjectValue("emailCodeSender", sender)
	}

	return sender, err
}

func initEmailCodeStore() (store code.CodeStore, err error) {
	if config.Setting.Enabled {
		if config.Setting.Store.Redis != nil && config.Setting.Store.Redis.Enabled {
			store, err = code.CreateRedisCodeStore(config.Setting.Store)
		} else if config.Setting.Store.Cache != nil && config.Setting.Store.Cache.Enabled {
			store, err = code.CreateCacheCodeStore(config.Setting.Store)
		} else {
			store, err = code.CreateMemoryCodeStore(config.Setting.Store)
		}
		if store != nil {
			inject.InjectValue("emailCodeStore", store)
		}
	}
	return store, err
}

func initEmailCodeValidator() {
	if config.Setting.Enabled {
		inject.InjectValue("emailCodeValidator", &EmailCodeValidator{})
	}
}

func initEmailCodeController() {
	if config.Setting.Enabled {
		inject.InjectValue("emailCodeController", &EmailCodeController{})
	}
}
