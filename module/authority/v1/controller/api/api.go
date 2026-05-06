package api

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources = &controller.Controllers{
	Base: "/auth/v1",
	Controllers: []controller.Controller{
		authorityController,
	},
}
