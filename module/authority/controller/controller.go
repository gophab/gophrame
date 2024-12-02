package controller

import (
	"github.com/gophab/gophrame/controller"

	"github.com/gophab/gophrame/core/logger"

	"github.com/gophab/gophrame/module/authority/controller/api"
	"github.com/gophab/gophrame/module/authority/controller/mapi"
	"github.com/gophab/gophrame/module/authority/controller/openapi"
)

func init() {
	logger.Info("Intitializing Module Authority Controllers...")
	controller.AddSchemaControllers("/api", api.Resources)
	controller.AddSchemaControllers("/mapi", mapi.Resources)
	// controller.AddSchemaControllers("/openapi/public", openapi.PublicResources)
	controller.AddSchemaControllers("/openapi", openapi.UserResources)
	controller.AddSchemaControllers("/openapi/admin", openapi.AdminResources)
}
