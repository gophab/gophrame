package controller

import (
	"github.com/gophab/gophrame/controller"

	"github.com/gophab/gophrame/core/logger"

	"github.com/gophab/gophrame/module/common/controller/api"
	"github.com/gophab/gophrame/module/common/controller/mapi"
	"github.com/gophab/gophrame/module/common/controller/openapi"
)

func init() {
	logger.Info("Intitializing Module Common Controllers...")
	controller.AddSchemaControllers("/api", api.Resources)
	controller.AddSchemaControllers("/mapi", mapi.Resources)
	controller.AddSchemaControllers("/openapi", openapi.UserResources)
	controller.AddSchemaControllers("/openapi/public", openapi.PublicResources)
	controller.AddSchemaControllers("/openapi/admin", openapi.AdminResources)

}
