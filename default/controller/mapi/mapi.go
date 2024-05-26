package mapi

import (
	"github.com/gophab/gophrame/default/controller/mapi/auth"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/permission"
	"github.com/gophab/gophrame/core/security"

	"github.com/gin-gonic/gin"
)

var Resources = &controller.Controllers{
	Base: "/mapi",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),      // oauth2 验证
		permission.NeedSystemUser(),       // 需要系统用户
		permission.CheckUserPermissions(), // 权限验证
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
