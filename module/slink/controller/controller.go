package controller

import (
	"github.com/gophab/gophrame/controller"
	"github.com/gophab/gophrame/core/starter"

	"github.com/gophab/gophrame/module/slink/controller/mapi"
	"github.com/gophab/gophrame/module/slink/controller/public"
)

func init() {
	controller.AddSchemaControllers("/mapi", mapi.Resources)

	starter.RegisterStarter(public.Start)
}
