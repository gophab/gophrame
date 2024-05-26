package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/util"
)

var Resources = &Controllers{
	Controllers: []Controller{},
}

func AddController(c Controller) {
	if b, _ := util.Contains(Resources.Controllers, c); !b {
		Resources.Controllers = append(Resources.Controllers, c)
	}
}

func AddControllers(cs ...Controller) {
	for _, c := range cs {
		AddController(c)
	}
}

func InitRouter(engine *gin.Engine) {
	Resources.InitRouter(engine.Group("/"))
}
