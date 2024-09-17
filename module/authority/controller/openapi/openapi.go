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
		authorityOpenController,
	},
}

var AdminResources = &controller.Controllers{
	Base: "/auth",
	Controllers: []controller.Controller{
		adminAuthorityOpenController,
	},
}
