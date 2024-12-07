package captcha

import (
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/captcha/config"
	"github.com/gophab/gophrame/core/form"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/webservice/response"
)

type CaptchaForm struct {
	CaptchaId    string `form:"captcha_id" json:"captcha_id"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value"`
	CaptchaKeep  bool   `form:"captcha_keep" json:"captcha_keep"`
}

func HandleCaptchaVerify(force bool) gin.HandlerFunc {
	captchaIdKey := config.Setting.CaptchaId
	captchaValueKey := config.Setting.CaptchaValue

	if !config.Setting.Enabled {
		return func(context *gin.Context) {
			context.Next()
		}
	}

	return func(context *gin.Context) {
		captchaId := context.Param(captchaIdKey)
		value := context.Param(captchaValueKey)
		keep := context.Param("captcha_keep") == "true"

		if captchaId == "" || value == "" {
			var data CaptchaForm
			if err := form.ShouldBind(context, &data); err == nil {
				captchaId = data.CaptchaId
				value = data.CaptchaValue
				keep = data.CaptchaKeep
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
					if strings.HasPrefix(seg, "captcha_keep=") {
						keep, _ = strconv.ParseBool(strings.TrimPrefix(seg, "captcha_keep="))
					}
				}
			}
		}

		if captchaId == "" || value == "" {
			if force {
				response.FailMessage(context, CaptchaCheckParamsInvalidCode, CaptchaCheckParamsInvalidMsg)
				return
			} else {
				context.Set("captcha", false)
				context.Next()
				return
			}
		}

		if captcha.Verify(captchaId, []byte(value)) {
			if keep {
				// captcha.Store.Set(captchaId, []byte(value))
				captchaService.Store.CreateCode(captchaId, "", value)
			}
			context.Set("captcha", true)
			context.Next()
		} else if force {
			response.FailMessage(context, CaptchaCheckFailCode, CaptchaCheckFailMsg)
		} else {
			context.Set("captcha", false)
			context.Next()
		}
	}
}
