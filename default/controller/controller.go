package controller

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/logger"

	"github.com/gophab/gophrame/default/controller/api"
	"github.com/gophab/gophrame/default/controller/mapi"
	"github.com/gophab/gophrame/default/controller/openapi"
)

func init() {
	logger.Info("Intitializing Module default Controllers...")
	controller.AddControllers(api.Resources, mapi.Resources, openapi.Resources)
}
