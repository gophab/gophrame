package swagger

import (
	"github.com/gophab/gophrame/core/engine"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Starting Core Swagger: ...", global.Swagger)
	if global.Swagger {
		// Swagger
		engine.Get().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // API 注释
	}
}
