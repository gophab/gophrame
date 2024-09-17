package mapi

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources = &controller.Controllers{
	Controllers: []controller.Controller{
		tenantMController,
		userMController,
		socialUserMController,
		organizationMController,
		organizationUserMController,
		roleMController,
	},
}
