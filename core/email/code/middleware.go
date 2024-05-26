package code

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/webservice/response"
)

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
			verificationCode := context.Request.Header.Get("X-Verification-Code")
			if verificationCode != "" {
				segs := strings.Split(verificationCode, ";")
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
				response.FailMessage(context, EmailCodeCheckParamsInvalidCode, EmailCodeCheckParamsInvalidMsg)
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
			response.FailMessage(context, EmailCodeCheckFailCode, EmailCodeCheckFailMsg)
		}
	}
}
