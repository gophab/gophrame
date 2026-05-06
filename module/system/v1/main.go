package system

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"
	"github.com/gophab/gophrame/core/starter"

	_ "github.com/gophab/gophrame/module/system/v1/controller"
	_ "github.com/gophab/gophrame/module/system/v1/security"
	_ "github.com/gophab/gophrame/module/system/v1/service"
)

const (
	MODULE_ID = 1
)

var _module = &module.Module{
	Name:        "SystemV1",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name, "v1")
	module.RegisterModule(_module)

	starter.RegisterStarter(Start)
	// 1. 加载Config...
}

func Start() {
	//
}
