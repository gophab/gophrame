package api

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/security"

	"github.com/wjshen/gophrame/controller/api/auth"

	"github.com/gin-gonic/gin"
)

var Resources *controller.Controllers = &controller.Controllers{
	Base: "/api",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),    // oauth2 验证
		security.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{
		userController,
		roleController,
		organizationController,
		organizationUserController,
		socialUserController,
		auth.Resources,
	},
}

func init() {
	inject.InjectValue("apiResources", Resources)
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
