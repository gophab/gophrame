package captcha

import (
	"bytes"
	"net/http"
	"path"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/captcha/config"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
)

type CaptchaController struct {
	controller.ResourceController
	CaptchaService *CaptchaService `inject:"captchaService"`
}

// 生成验证码
func (c *CaptchaController) Generate(context *gin.Context) {
	var gtype = request.Param(context, "type").DefaultString("image")

	if captcha, err := c.CaptchaService.Generate(gtype); err == nil {
		response.Success(context, captcha)
	} else {
		response.SystemErrorMessage(context, 500, err.Error())
	}
}

// 获取验证码图像
func (c *CaptchaController) GetImg(context *gin.Context) {
	captchaIdKey := config.Setting.CaptchaId
	captchaId := context.Param(captchaIdKey)
	_, file := path.Split(context.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || captchaId == "" {
		response.FailMessage(context, CaptchaGetParamsInvalidCode, CaptchaGetParamsInvalidMsg)
		return
	}

	if context.Query("reload") != "" {
		captcha.Reload(id)
	}

	context.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Header("Pragma", "no-cache")
	context.Header("Expires", "0")

	var vBytes bytes.Buffer
	if ext == ".png" {
		context.Header("Content-Type", "image/png")
		// 设置实际业务需要的验证码图片尺寸（宽 X 高），captcha.StdWidth, captcha.StdHeight 为默认值，请自行修改为具体数字即可
		_ = captcha.WriteImage(&vBytes, id, captcha.StdWidth, captcha.StdHeight)
		http.ServeContent(context.Writer, context.Request, id+ext, time.Time{}, bytes.NewReader(vBytes.Bytes()))
	}
}

// 校验验证码
func (c *CaptchaController) CheckCode(context *gin.Context) {
	captchaIdKey := config.Setting.CaptchaId
	captchaValueKey := config.Setting.CaptchaValue

	captchaId := context.Param(captchaIdKey)
	value := context.Param(captchaValueKey)

	if captchaId == "" {
		captchaId = context.Request.Header.Get(captchaIdKey)
	}

	if value == "" {
		value = context.Request.Header.Get(captchaValueKey)
	}

	if captchaId == "" || value == "" {
		response.FailMessage(context, CaptchaCheckParamsInvalidCode, CaptchaCheckParamsInvalidMsg)
		return
	}

	if c.CaptchaService.Verify(captchaId, value) {
		response.Success(context, CaptchaCheckOkMsg)
	} else {
		response.FailMessage(context, CaptchaCheckFailCode, CaptchaCheckFailMsg)
	}
}

func (c *CaptchaController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	captcha := g.Group("openapi/captcha")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		captcha.GET("", c.Generate)                             //  获取验证码ID
		captcha.GET("/:captcha_id", c.GetImg)                   // 获取图像地址
		captcha.GET("/:captcha_id/:captcha_value", c.CheckCode) // 校验验证码
	}
	return captcha
}
