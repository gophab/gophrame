package auth

import (
	"github.com/wjshen/gophrame/core/controller"

	"github.com/gin-gonic/gin"
)

var Resources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		authorityOpenController,
		buttonOpenController,
		menuOpenController,
	},
}

func InitRouter(r *gin.RouterGroup) {
	Resources.InitRouter(r)
}
