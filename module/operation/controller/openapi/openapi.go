package openapi

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources = &controller.Controllers{
	Controllers: []controller.Controller{
		UserResources,
		AdminResources,
		// PublicResources,
	},
}

var UserResources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		// buttonOpenController,
		// menuOpenController,
		moduleOpenController,
	},
}

var AdminResources = &controller.Controllers{
	Base:        "/auth",
	Controllers: []controller.Controller{},
}
