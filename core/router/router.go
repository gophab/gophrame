package router

import (
	"net/http"

	"github.com/gophab/gophrame/core/engine"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	_ "github.com/gophab/gophrame/core/swagger"

	"github.com/gin-gonic/gin"
)

func init() {
}

func Root() *gin.Engine {
	return engine.Get()
}

func Init() {
	logger.Debug("Initializing Core Router...")

	engine.Init(global.Debug)

	// 处理
	Root().GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "OK")
	})

	//处理静态资源
	Root().Static("/public", "./public") //  定义静态资源路由与实际目录映射关系
}

func Start() {
	logger.Info("Starting Core Router...")
}
