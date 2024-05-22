package code

import (
	"strings"

	"github.com/wjshen/gophrame/core/captcha"
	"github.com/wjshen/gophrame/core/code"
	"github.com/wjshen/gophrame/core/form"
	"github.com/wjshen/gophrame/core/sms"
	"github.com/wjshen/gophrame/core/sms/config"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/wjshen/gophrame/errors"

	"github.com/gin-gonic/gin"
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
	params["product"] = config.Setting.Product
	params["signature"] = config.Setting.Signature

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

type PhoneForm struct {
	Phone string `form:"phone"`
	Code  string `form:"code"`
	Scene string `form:"scene"`
}

func (v *SmsCodeValidator) HandleSmsCodeVerify(force bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		phone := context.Param("phone")
		scene := context.Param("scene")
		value := context.Param("code")

		if phone == "" || value == "" || scene == "" {
			var data PhoneForm
			if form.ShouldBind(context, &data) == nil {
				phone = data.Phone
				value = data.Code
				scene = data.Scene
			}
		}

		if phone == "" || value == "" || scene == "" {
			verificationCode := context.Request.Header.Get("X-Verification-Code")
			if verificationCode != "" {
				segs := strings.Split(verificationCode, ";")
				for _, seg := range segs {
					seg = strings.TrimSpace(seg)
					if strings.HasPrefix(seg, "phone=") {
						phone = strings.TrimPrefix(seg, "phone=")
					}
					if strings.HasPrefix(seg, "scene=") {
						scene = strings.TrimPrefix(seg, "scene=")
					}
					if strings.HasPrefix(seg, "code=") {
						value = strings.TrimPrefix(seg, "code=")
					}
				}
			}
		}

		if phone == "" || value == "" {
			if force {
				response.Fail(context, SmsCodeCheckParamsInvalidCode, SmsCodeCheckParamsInvalidMsg)
				return
			} else {
				context.AddParam("smscode", "false")
				context.Next()
			}
		}

		if b := v.CheckCode(v, phone, scene, value); b {
			context.AddParam("smscode", "true")
			context.Next()
		} else {
			response.Fail(context, SmsCodeCheckFailCode, SmsCodeCheckFailMsg)
		}
	}
}

/**
 * Controller
 */
type SmsCodeController struct {
	SmsCodeValidator *SmsCodeValidator `inject:"smsCodeValidator"`
}

func (s *SmsCodeController) GenerateCode(c *gin.Context) {
	phone, err := request.Param(c, "phone").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	scene := request.Param(c, "scene").DefaultString("register-pin")
	force := request.Param(c, "force").DefaultBool(false)

	_, b := s.SmsCodeValidator.GetVerificationCode(s.SmsCodeValidator, phone, scene)
	if !b || force {
		_, err = s.SmsCodeValidator.GenerateCode(s.SmsCodeValidator, phone, scene)
		if err != nil {
			response.FailMessage(c, 400, err.Error())
			return
		}
	} /* else 验证码仍旧有效，忽略 */

	response.OK(c, "验证码已发送")
}

func (s *SmsCodeController) CheckCode(c *gin.Context) {
	phone, err := request.Param(c, "phone").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	code, err := request.Param(c, "code").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	scene := request.Param(c, "scene").DefaultString("register-pin")

	response.OK(c, s.SmsCodeValidator.CheckCode(s.SmsCodeValidator, phone, scene, code))
}

/**
 * 处理WEB验证码的API路由
 */
func (s *SmsCodeController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		verifyCode := g.Group("openapi/sms")
		{
			verifyCode.Use(captcha.HandleCaptchaVerify(false))

			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			verifyCode.GET("/code", s.GenerateCode)           // 发送验证码
			verifyCode.GET("/code/:phone/:code", s.CheckCode) // 校验验证码
		}
	}
}
