package code

import (
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/captcha"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

/**
 * Controller
 */
type SmsCodeController struct {
	controller.ResourceController
	SmsCodeValidator *SmsCodeValidator `inject:"smsCodeValidator"`
}

func (s *SmsCodeController) GenerateCode(c *gin.Context) {
	phone, err := request.Param(c, "phone").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	scene := request.Param(c, "scene").DefaultString("register")
	force := request.Param(c, "force").DefaultBool(false)

	_, b := s.SmsCodeValidator.GetVerificationCode(s.SmsCodeValidator, phone, scene)
	if !b || force {
		_, err = s.SmsCodeValidator.GenerateCode(s.SmsCodeValidator, phone, scene)
		if err != nil {
			response.FailMessage(c, 400, err.Error())
			return
		}
	} /* else 验证码仍旧有效，忽略 */

	response.Success(c, "验证码已发送")
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

	response.Success(c, s.SmsCodeValidator.CheckCode(s.SmsCodeValidator, phone, scene, code))
}

/**
 * 处理WEB验证码的API路由
 */
func (s *SmsCodeController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	smsRoot := g.Group("openapi/sms")
	{
		smsRoot.Use(captcha.HandleCaptchaVerify(false))

		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		smsRoot.GET("/code", s.GenerateCode)           // 发送验证码
		smsRoot.GET("/code/:phone/:code", s.CheckCode) // 校验验证码
	}

	return smsRoot
}
