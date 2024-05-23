package controller

import (
	"github.com/wjshen/gophrame/core/controller"

	"github.com/gin-gonic/gin"
)

var Resources = &controller.Controllers{
	Controllers: []controller.Controller{
		// security APIs
		securityController,
	},
}

func AddController(c controller.Controller) {
	Resources.Controllers = append(Resources.Controllers, c)
}

func AddControllers(cs ...controller.Controller) {
	Resources.Controllers = append(Resources.Controllers, cs...)
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
