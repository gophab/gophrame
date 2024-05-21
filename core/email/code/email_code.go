package code

import (
	"strings"

	"github.com/wjshen/gophrame/core/captcha"
	"github.com/wjshen/gophrame/core/code"
	"github.com/wjshen/gophrame/core/email"
	"github.com/wjshen/gophrame/core/email/config"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"

	"github.com/gin-gonic/gin"
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

	return s.EmailSender.SendTemplateEmail(dest, scene, params)
}

func (v *EmailCodeValidator) GetSender() code.CodeSender {
	if v.Sender == nil {
		if config.Setting.Enabled {
		} else {
			return v.Validator.GetSender()
		}
	}
	return v.Sender
}

func (v *EmailCodeValidator) GetStore() code.CodeStore {
	if v.Store == nil {
		if config.Setting.Enabled {
			if config.Setting.Store.Redis != nil && config.Setting.Store.Redis.Enabled {
				v.Store, _ = code.CreateRedisCodeStore(config.Setting.Store)
			} else if config.Setting.Store.Cache != nil && config.Setting.Store.Cache.Enabled {
				v.Store, _ = code.CreateCacheCodeStore(config.Setting.Store)
			} else {
				v.Store, _ = code.CreateMemoryCodeStore(config.Setting.Store)
			}
		} else {
			return v.Validator.GetStore()
		}
	}
	return v.Store
}

type EmailForm struct {
	Email string `form:"phone"`
	Code  string `form:"code"`
}

func (v *EmailCodeValidator) HandleSmsCodeVerify(force bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		email := context.Param("email")
		scene := context.Param("scene")
		value := context.Param("code")

		if email == "" || value == "" || scene == "" {
			captcha := context.Request.Header.Get("X-Verification-Code")
			if captcha != "" {
				segs := strings.Split(captcha, ";")
				for _, seg := range segs {
					seg = strings.TrimSpace(seg)
					if strings.HasPrefix(seg, "email=") {
						email = strings.TrimPrefix(seg, "email=")
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

		if email == "" || value == "" || scene == "" {
			if force {
				response.Fail(context, EmailCodeCheckParamsInvalidCode, EmailCodeCheckParamsInvalidMsg)
				return
			} else {
				context.AddParam("emailcode", "false")
				context.Next()
			}
		}

		if b := v.CheckCode(v, email, scene, value); b {
			context.AddParam("emailcode", "true")
			context.Next()
		} else {
			response.Fail(context, EmailCodeCheckFailCode, EmailCodeCheckFailMsg)
		}
	}
}

type EmailCodeController struct {
	EmailCodeValidator *code.Validator `inject:"emailCodeValidator"`
}

func (e *EmailCodeController) GenerateCode(c *gin.Context) {
	email, err := request.Param(c, "email").MustString()
	scene := request.Param(c, "scene").DefaultString("register-pin")
	force := request.Param(c, "force").DefaultBool(false)

	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	_, b := e.EmailCodeValidator.GetVerificationCode(e.EmailCodeValidator, email, scene)
	if !b || force {
		_, err = e.EmailCodeValidator.GenerateCode(e.EmailCodeValidator, email, scene)
		if err != nil {
			response.FailMessage(c, 400, err.Error())
			return
		}
	} /* else 验证码仍旧有效，忽略 */

	response.OK(c, "验证码已发送")
}

func (e *EmailCodeController) CheckCode(c *gin.Context) {
	email, err := request.Param(c, "email").MustString()
	scene := request.Param(c, "scene").DefaultString("register-pin")
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	code, err := request.Param(c, "code").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	response.OK(c, e.EmailCodeValidator.CheckCode(e.EmailCodeValidator, email, scene, code))
}

/**
 * 处理WEB验证码的API路由
 */
func (e *EmailCodeController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		verifyCode := g.Group("email")
		{
			verifyCode.Use(captcha.HandleCaptchaVerify(false))

			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			verifyCode.GET("/code", e.GenerateCode)            // 发送验证码
			verifyCode.GET("/check/:email/:code", e.CheckCode) // 校验验证码
		}
	}
}
