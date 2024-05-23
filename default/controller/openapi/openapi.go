package openapi

import (
	"github.com/wjshen/gophrame/default/controller/openapi/auth"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/security"

	"github.com/gin-gonic/gin"
)

var Resources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		UserResources,
		AdminResources,
		PublicResources,
	},
}

var PublicResources *controller.Controllers = &controller.Controllers{
	Base: "/openapi/public",
	Handlers: []gin.HandlerFunc{
		security.CheckTokenVerify(),     // oauth2 验证
		security.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{
		publicSystemOptionOpenController,
	},
}

var UserResources *controller.Controllers = &controller.Controllers{
	Base: "/openapi/user",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),    // oauth2 验证
		security.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{
		userOpenController,
		roleOpenController,
		organizationOpenController,
		organizationUserOpenController,
		userOptionOpenController,
		socialUserOpenController,
		tenantOptionOpenController,
		auth.Resources,
	},
}

var AdminResources *controller.Controllers = &controller.Controllers{
	Base: "/openapi/admin",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),    // oauth2 验证
		security.CheckUserPermissions(), // 权限验证
		security.NeedAdmin(),
	},
	Controllers: []controller.Controller{
		adminTenantOptionOpenController,
	},
}

func init() {
	inject.InjectValue("userResources", UserResources)
	inject.InjectValue("adminResources", AdminResources)
	inject.InjectValue("publicResources", PublicResources)
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
