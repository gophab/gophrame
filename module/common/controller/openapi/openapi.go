package openapi

import (
	"github.com/gophab/gophrame/core/controller"
)

var Resources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		UserResources,
		AdminResources,
		PublicResources,
	},
}

var PublicResources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		publicSystemOptionOpenController,
	},
}

var UserResources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		eventOpenController,
		messageOpenController,
		tenantOptionOpenController,
		userOptionOpenController,
	},
}

var AdminResources *controller.Controllers = &controller.Controllers{
	Controllers: []controller.Controller{
		adminEventOpenController,
		adminMessageOpenController,
		adminTenantOptionOpenController,
	},
}
