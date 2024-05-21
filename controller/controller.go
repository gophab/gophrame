package controller

import (
	"github.com/wjshen/gophrame/controller/openapi"
	"github.com/wjshen/gophrame/core/controller"

	"github.com/gin-gonic/gin"
)

var Resources = &controller.Controllers{
	Controllers: []controller.Controller{
		// security APIs
		securityController,

		// open APIs
		openapi.Resources,

		// internal management APIs
		//mapi.Resources,

		// internal service APIs
		//api.Resources,
	},
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
