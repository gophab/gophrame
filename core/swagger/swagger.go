package swagger

import (
	"github.com/gophab/gophrame/core/engine"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
