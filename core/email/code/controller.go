package code

import (
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/captcha"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type EmailCodeController struct {
	controller.ResourceController
	EmailCodeValidator *EmailCodeValidator `inject:"emailCodeValidator"`
}

func (e *EmailCodeController) GenerateCode(c *gin.Context) {
	email, err := request.Param(c, "email").MustString()
	scene := request.Param(c, "scene").DefaultString("register")
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

	response.Success(c, "验证码已发送")
}

func (e *EmailCodeController) CheckCode(c *gin.Context) {
	email, err := request.Param(c, "email").MustString()
	scene := request.Param(c, "scene").DefaultString("register")
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	code, err := request.Param(c, "code").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	response.Success(c, e.EmailCodeValidator.CheckCode(e.EmailCodeValidator, email, scene, code))
}

/**
 * 处理WEB验证码的API路由
 */
func (e *EmailCodeController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	email := g.Group("openapi/email")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		email.GET("/code", captcha.HandleCaptchaVerify(true), e.GenerateCode) // 发送验证码
		email.GET("/check/:email/:code", e.CheckCode)                         // 校验验证码
	}

	return email
}
