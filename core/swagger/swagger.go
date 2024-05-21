package swagger

import (
	"github.com/wjshen/gophrame/core/engine"
	_ "github.com/wjshen/gophrame/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	// Swagger
	engine.Get().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // API 注释
}
