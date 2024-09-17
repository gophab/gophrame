package openapi

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		UserResources,
		AdminResources,
	},
}

var UserResources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		organizationOpenController,
		organizationUserOpenController,
		userOpenController,
		socialUserOpenController,
	},
}

var AdminResources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		adminUserOpenController,
		adminOrganizationOpenController,
		adminTenantOpenController,
		adminRoleOpenController,
	},
}
