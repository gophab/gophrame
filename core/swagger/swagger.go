package swagger

import (
	"github.com/wjshen/gophrame/core/engine"

	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	// Swagger
	engine.Get().GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler)) // API 注释
}
