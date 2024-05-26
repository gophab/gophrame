package auth

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		authorityOpenController,
		buttonOpenController,
		menuOpenController,
	},
}
