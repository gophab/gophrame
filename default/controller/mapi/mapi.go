package mapi

import (
	"github.com/wjshen/gophrame/default/controller/mapi/auth"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/security"

	"github.com/gin-gonic/gin"
)

var Resources = &controller.Controllers{
	Base: "/mapi",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),    // oauth2 验证
		security.CheckUserPermissions(), // 权限验证
		security.NeedSystemUser(),
	},
	Controllers: []controller.Controller{
		userMController,
		roleMController,
		socialUserMController,
		organizationMController,
		organizationUserMController,
		systemOptionMController,
		tenantOptionMController,
		auth.Resources,
	},
}

func init() {
	inject.InjectValue("Resources", Resources)
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
