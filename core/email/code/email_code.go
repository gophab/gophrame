package code

import (
	"github.com/gophab/gophrame/core/code"
	"github.com/gophab/gophrame/core/email"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/core/email/code/config"
)

const (
	//验证码
	EmailCodeGetParamsInvalidMsg    string = "获取验证码：提交的验证码参数无效,请检查验证码ID以及文件名后缀是否完整"
	EmailCodeGetParamsInvalidCode   int    = -400350
	EmailCodeCheckParamsInvalidMsg  string = "校验验证码：提交的参数无效，请检查 【验证码ID、验证码值】 提交时的键名是否与配置项一致"
	EmailCodeCheckParamsInvalidCode int    = -400351
	EmailCodeCheckOkMsg             string = "验证码校验通过"
	EmailCodeCheckFailCode          int    = -400355
	EmailCodeCheckFailMsg           string = "验证码校验失败"
)

type EmailCodeValidator struct {
	code.Validator
	Sender code.CodeSender `inject:"emailCodeSender"`
	Store  code.CodeStore  `inject:"emailCodeStore"`
}

type EmailCodeSender struct {
	EmailSender email.EmailSender `inject:"emailSender"`
}

func (s *EmailCodeSender) SendVerificationCode(dest string, scene string, code string) error {
	params := make(map[string]string)
	params["code"] = code

	t := scene

	// TODO: 通用模板处理
	// t := template.GetTemplate("email:" + scene)

	return s.EmailSender.SendTemplateEmail(dest, t, params)
}

func (v *EmailCodeValidator) GetSender() code.CodeSender {
	if v.Sender == nil {
		if config.Setting.Enabled {
			sender := &EmailCodeSender{}
			inject.InjectValue("emailCodeSender", sender)
			return sender
		} else {
			return v.Validator.GetSender()
		}
	}
	return v.Sender
}

func (v *EmailCodeValidator) GetStore() code.CodeStore {
	if v.Store == nil {
		if config.Setting.Enabled {
			if config.Setting.Redis != nil && config.Setting.Redis.Enabled {
				v.Store, _ = code.CreateRedisCodeStore(config.Setting)
			} else if config.Setting.Cache != nil && config.Setting.Cache.Enabled {
				v.Store, _ = code.CreateCacheCodeStore(config.Setting)
			} else {
				v.Store, _ = code.CreateMemoryCodeStore(config.Setting)
			}
		} else {
			return v.Validator.GetStore()
		}
	}
	return v.Store
}
