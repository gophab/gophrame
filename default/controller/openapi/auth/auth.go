package auth

import (
	"github.com/gophab/gophrame/core/controller"
)

var UserResources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		authorityOpenController,
		buttonOpenController,
		menuOpenController,
		moduleOpenController,
	},
}

var AdminResources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		adminAuthorityOpenController,
	},
}
