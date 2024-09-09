package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/context"
	"github.com/gophab/gophrame/core/webservice/request"
)

func SetGlobalContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context.SetContextValue("_current_context_", ctx)
		defer context.RemoveContextValue("_current_context_")

		ctx.Next()
	}
}

func EnableLocale() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetHeader("x-set-locale")
		if locale == "" {
			locale = request.Param(ctx, "_LOCALE_").DefaultString("")
		}

		if locale != "" {
			context.SetContextValue("_LOCALE_", locale)
		}
		defer context.RemoveContextValue("_LOCALE_")

		ctx.Next()
	}
}
