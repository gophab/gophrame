package controller

import (
	"github.com/gophab/gophrame/controller"
	"github.com/gophab/gophrame/core/logger"

	"github.com/gophab/gophrame/module/operation/controller/api"
	"github.com/gophab/gophrame/module/operation/controller/mapi"
	"github.com/gophab/gophrame/module/operation/controller/openapi"
)

func init() {
	logger.Info("Intitializing Module Operation Controllers...")
	controller.AddSchemaControllers("/api", api.Resources)
	controller.AddSchemaControllers("/mapi", mapi.Resources)
	// controller.AddSchemaControllers("/openapi/public", openapi.PublicResources)
	controller.AddSchemaControllers("/openapi/user", openapi.UserResources)
	controller.AddSchemaControllers("/openapi/admin", openapi.AdminResources)

}
