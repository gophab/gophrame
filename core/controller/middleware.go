package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/context"
)

func SetGlobalContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context.SetContextValue("_current_context_", ctx)
		defer context.RemoveContextValue("_current_context_")

		ctx.Next()
	}
}
