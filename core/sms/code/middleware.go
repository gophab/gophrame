package code

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/form"
	"github.com/gophab/gophrame/core/webservice/response"
)

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
				response.FailMessage(context, SmsCodeCheckParamsInvalidCode, SmsCodeCheckParamsInvalidMsg)
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
			response.FailMessage(context, SmsCodeCheckFailCode, SmsCodeCheckFailMsg)
		}
	}
}
