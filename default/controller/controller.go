package controller

import (
	"github.com/wjshen/gophrame/controller"
	"github.com/wjshen/gophrame/default/controller/api"
	"github.com/wjshen/gophrame/default/controller/mapi"
	"github.com/wjshen/gophrame/default/controller/openapi"
)

func init() {
	controller.AddControllers(api.Resources, mapi.Resources, openapi.Resources)
}
