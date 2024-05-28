package captcha

import (
	"github.com/gophab/gophrame/core/captcha/config"
	"github.com/gophab/gophrame/core/code"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

const (
	//验证码
	CaptchaGetParamsInvalidMsg    string = "获取验证码：提交的验证码参数无效,请检查验证码ID以及文件名后缀是否完整"
	CaptchaGetParamsInvalidCode   int    = -400350
	CaptchaCheckParamsInvalidMsg  string = "校验验证码：提交的参数无效，请检查 【验证码ID、验证码值】 提交时的键名是否与配置项一致"
	CaptchaCheckParamsInvalidCode int    = -400351
	CaptchaCheckOkMsg             string = "验证码校验通过"
	CaptchaCheckOkCode            int    = 200
	CaptchaCheckFailCode          int    = -400355
	CaptchaCheckFailMsg           string = "图形验证码校验失败"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing Captcha ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		var store code.CodeStore
		if config.Setting.Store.Enabled {
			if config.Setting.Store.Cache != nil && config.Setting.Store.Cache.Enabled {
				store, _ = code.CreateCacheCodeStore(config.Setting.Store)
			} else if config.Setting.Store.Redis != nil && config.Setting.Store.Redis.Enabled {
				store, _ = code.CreateRedisCodeStore(config.Setting.Store)
			} else {
				store, _ = code.CreateMemoryCodeStore(config.Setting.Store)
			}
		}

		service := &CaptchaService{
			Store: store,
		}
		inject.InjectValue("captchaService", service)

		service.Init()

		controller.AddController(&CaptchaController{
			CaptchaService: service,
		})
	}
}
