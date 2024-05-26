package code

import (
	"sync"

	"github.com/gophab/gophrame/core/code"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/email/code/config"
	"github.com/gophab/gophrame/core/inject"
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
		if config.Setting.Redis != nil && config.Setting.Redis.Enabled {
			store, err = code.CreateRedisCodeStore(config.Setting)
		} else if config.Setting.Cache != nil && config.Setting.Cache.Enabled {
			store, err = code.CreateCacheCodeStore(config.Setting)
		} else {
			store, err = code.CreateMemoryCodeStore(config.Setting)
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
		var emailCodeController = &EmailCodeController{}
		inject.InjectValue("emailCodeController", emailCodeController)
		controller.AddController(emailCodeController)
	}
}
