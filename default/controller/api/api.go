package api

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/security"

	"github.com/gophab/gophrame/default/controller/api/auth"

	"github.com/gin-gonic/gin"
)

var Resources *controller.Controllers = &controller.Controllers{
	Base: "/api",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(), // oauth2 验证
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
