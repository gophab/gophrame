package openapi

import (
	"github.com/gophab/gophrame/default/controller/openapi/auth"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/permission"
	"github.com/gophab/gophrame/core/security"

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
		security.CheckTokenVerify(), // oauth2 验证
	},
	Controllers: []controller.Controller{
		publicSystemOptionOpenController,
	},
}

var UserResources *controller.Controllers = &controller.Controllers{
	Base: "/openapi/user",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),      // oauth2 验证
		permission.CheckUserPermissions(), // 权限验证
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
		security.HandleTokenVerify(), // oauth2 验证
		permission.NeedAdmin(),
		permission.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{
		adminOrganizationOpenController,
		adminTenantOpenController,
		adminTenantOptionOpenController,
	},
}
