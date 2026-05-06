package mapi

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources = &controller.Controllers{
	Base: "/v1",
	Controllers: []controller.Controller{
		tenantMController,
		userMController,
		socialUserMController,
		organizationMController,
		organizationUserMController,
		roleMController,
	},
}
