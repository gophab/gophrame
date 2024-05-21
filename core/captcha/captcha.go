package captcha

import (
	"bytes"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/wjshen/gophrame/core/captcha/config"
	"github.com/wjshen/gophrame/core/form"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
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

type Captcha struct{}

// 生成验证码ID
func (c *Captcha) GenerateId(context *gin.Context) {
	// 设置验证码的数字长度（个数）
	var length = config.Setting.Length
	var captchaId, imgUrl, refresh, verify string

	captchaId = captcha.NewLen(length)
	imgUrl = "/openapi/captcha/" + captchaId + ".png"
	refresh = imgUrl + "?reload=1"
	verify = "/openapi/captcha/" + captchaId + "/{code}"

	response.Success(context, gin.H{
		"id":      captchaId,
		"imgUrl":  imgUrl,
		"refresh": refresh,
		"verify":  verify,
	})
}

// 获取验证码图像
func (c *Captcha) GetImg(context *gin.Context) {
	captchaIdKey := config.Setting.CaptchaId
	captchaId := context.Param(captchaIdKey)
	_, file := path.Split(context.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || captchaId == "" {
		response.Fail(context, CaptchaGetParamsInvalidCode, CaptchaGetParamsInvalidMsg)
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
func (c *Captcha) CheckCode(context *gin.Context) {
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
		response.Fail(context, CaptchaCheckParamsInvalidCode, CaptchaCheckParamsInvalidMsg)
		return
	}
	if captcha.VerifyString(captchaId, value) {
		response.Success(context, CaptchaCheckOkMsg)
	} else {
		response.Fail(context, CaptchaCheckFailCode, CaptchaCheckFailMsg)
	}
}

type CaptchaForm struct {
	CaptchaId    string `form:"captcha_id" json:"captcha_id"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value"`
}

func HandleCaptchaVerify(force bool) gin.HandlerFunc {
	captchaIdKey := config.Setting.CaptchaId
	captchaValueKey := config.Setting.CaptchaValue

	return func(context *gin.Context) {
		captchaId := context.Param(captchaIdKey)
		value := context.Param(captchaValueKey)

		if captchaId == "" || value == "" {
			var data CaptchaForm
			if err := form.ShouldBind(context, &data); err == nil {
				captchaId = data.CaptchaId
				value = data.CaptchaValue
			} else {
				logger.Warn("Captch parameters bind error: ", err.Error())
			}
		}

		if captchaId == "" || value == "" {
			captcha := context.Request.Header.Get("X-Verification-Code")
			if captcha != "" {
				segs := strings.Split(captcha, ";")
				for _, seg := range segs {
					seg = strings.TrimSpace(seg)
					if strings.HasPrefix(seg, captchaIdKey+"=") {
						captchaId = strings.TrimPrefix(seg, captchaIdKey+"=")
					}
					if strings.HasPrefix(seg, captchaValueKey+"=") {
						value = strings.TrimPrefix(seg, captchaValueKey+"=")
					}
				}
			}
		}

		if captchaId == "" || value == "" {
			if force {
				response.Fail(context, CaptchaCheckParamsInvalidCode, CaptchaCheckParamsInvalidMsg)
				return
			} else {
				context.AddParam("captcha", "false")
				context.Next()
				return
			}
		}

		if captcha.VerifyString(captchaId, value) {
			context.AddParam("captcha", "true")
			context.Next()
		} else {
			response.Fail(context, CaptchaCheckFailCode, CaptchaCheckFailMsg)
		}
	}
}

/**
 * 处理WEB验证码的API路由
 */
func InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		verifyCode := g.Group("openapi/captcha")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			verifyCode.GET("/", (&Captcha{}).GenerateId)                          //  获取验证码ID
			verifyCode.GET("/:captcha_id", (&Captcha{}).GetImg)                   // 获取图像地址
			verifyCode.GET("/:captcha_id/:captcha_value", (&Captcha{}).CheckCode) // 校验验证码
		}
	}
}
