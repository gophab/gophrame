package code

import (
	"github.com/gophab/gophrame/core/code"
	"github.com/gophab/gophrame/core/sms"
)

const (
	//验证码
	SmsCodeGetParamsInvalidMsg    string = "获取验证码：提交的验证码参数无效,请检查验证码ID以及文件名后缀是否完整"
	SmsCodeGetParamsInvalidCode   int    = -400350
	SmsCodeCheckParamsInvalidMsg  string = "校验验证码：提交的参数无效，请检查 【验证码ID、验证码值】 提交时的键名是否与配置项一致"
	SmsCodeCheckParamsInvalidCode int    = -400351
	SmsCodeCheckOkMsg             string = "验证码校验通过"
	SmsCodeCheckFailCode          int    = -400355
	SmsCodeCheckFailMsg           string = "验证码校验失败"
)

type SmsCodeSender struct {
	sms.SmsSender `inject:"smsSender"`
}

func (s *SmsCodeSender) SendVerificationCode(dest string, scene string, code string) error {
	params := map[string]string{}
	params["code"] = code
	// params["product"] = config.Setting.Product
	// params["signature"] = config.Setting.Signature

	return s.SendTemplateMessage(dest, scene, params)
}

/**
 * Validator
 */
type SmsCodeValidator struct {
	code.Validator
	Sender code.CodeSender `inject:"smsCodeSender"`
	Store  code.CodeStore  `inject:"smsCodeStore"`
}

func (v *SmsCodeValidator) GetSender() code.CodeSender {
	if v.Sender != nil {
		return v.Sender
	}
	return v.Validator.GetSender()
}

func (v *SmsCodeValidator) GetStore() code.CodeStore {
	if v.Store != nil {
		return v.Store
	}
	return v.Validator.GetStore()
}
