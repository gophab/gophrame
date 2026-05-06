package api

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources *controller.Controllers = &controller.Controllers{
	Base: "/v1",
	Controllers: []controller.Controller{
		userController,
		roleController,
		organizationController,
		organizationUserController,
		socialUserController,
	},
}
