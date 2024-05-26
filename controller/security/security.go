package security

import (
	"github.com/gophab/gophrame/core/captcha"
	"github.com/gophab/gophrame/core/controller"
	EmailCode "github.com/gophab/gophrame/core/email/code"
	SmsCode "github.com/gophab/gophrame/core/sms/code"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"
	"github.com/gophab/gophrame/service"
	"github.com/gophab/gophrame/service/dto"

	"github.com/gin-gonic/gin"
)

type SecurityController struct {
	controller.ResourceController
	MobileValidator   *SmsCode.SmsCodeValidator     `inject:"smsCodeValidator"`
	EmailValidator    *EmailCode.EmailCodeValidator `inject:"emailCodeValidator"`
	UserService       service.UserService           `inject:"userService"`
	InviteCodeService service.InviteCodeService     `inject:"inviteCodeService"`
}

func (c *SecurityController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	g.POST("/openapi/register", captcha.HandleCaptchaVerify(false), c.Register)      // 注册
	g.PUT("/openapi/register", captcha.HandleCaptchaVerify(false), c.ChangeRegister) // 注册
	return g
}

type RegisterForm struct {
	Mode       string `json:"mode" form:"mode"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	Password2  string `json:"password2,omitempty" form:"password2"`
	InviteCode string `json:"invite_code" form:"invite_code"`
}

// @Summary 注册新用户 （由用户主动发起）
// @Router /register [POST]
func (u *SecurityController) Register(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		response.FailMessage(c, errors.INVALID_PARAMS, err.Error())
		return
	}

	user := &dto.User{}

	switch form.Mode {
	case "password":
		{
			if form.Username == "" || form.Password == "" {
				response.FailMessage(c, errors.INVALID_PARAMS, "用户名或密码不能为空")
				return
			}

			if form.Password != form.Password2 {
				response.FailMessage(c, errors.INVALID_PARAMS, "两次输入密码不一致")
				return
			}

			user.Login = &form.Username
			user.Password = &form.Password

			if user.Name == nil {
				user.Name = &form.Username
			}
		}
	case "mobile":
		{
			if form.Username == "" || form.Password == "" {
				response.FailMessage(c, errors.INVALID_PARAMS, "手机号或验证码不能为空")
				return
			}

			if u.MobileValidator == nil {
				response.FailMessage(c, errors.INVALID_PARAMS, "不支持手机验证码登录")
				return
			}

			if !u.MobileValidator.CheckCode(u.MobileValidator, form.Username, "register-pin", form.Password) {
				response.FailMessage(c, errors.INVALID_PARAMS, "验证码不一致")
				return
			}

			user.Mobile = &form.Username
			if user.Name == nil {
				user.Name = &form.Username
			}
		}
	case "email":
		{
			if form.Username == "" || form.Password == "" {
				response.FailMessage(c, errors.INVALID_PARAMS, "Email或验证码不能为空")
				return
			}

			if u.EmailValidator == nil {
				response.FailMessage(c, errors.INVALID_PARAMS, "不支持Email验证码登录")
				return
			}

			if !u.EmailValidator.CheckCode(u.EmailValidator, form.Username, "register-pin", form.Password) {
				response.FailMessage(c, errors.INVALID_PARAMS, "验证码不一致")
				return
			}

			user.Email = &form.Username
			if user.Name == nil {
				user.Name = &form.Username
			}
		}
	}

	if form.InviteCode != "" {
		// 验证邀请码
		if iv, err := u.InviteCodeService.FindByInviteCode(form.InviteCode); err != nil {
			response.SystemErrorMessage(c, 400, err.Error())
			return
		} else if iv == nil {
			response.FailMessage(c, errors.INVALID_PARAMS, "无效的邀请码或邀请码已过期")
			return
		} else {
			user.InviteCode = &iv.InviteCode
			user.InviterId = &iv.UserId
		}
	} else {
		response.FailMessage(c, errors.INVALID_PARAMS, "邀请码不能为空")
		return
	}

	if res, err := u.UserService.CreateUser(user); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	} else {
		response.Success(c, res)
	}
}

/**
 * PUT /register
 */
func (o *SecurityController) ChangeRegister(c *gin.Context) {
	// 1. Change Login => 修改账号
	// 2. Change Mobile => 修改手机号
	// 3. Change Email => 修改Email
}
