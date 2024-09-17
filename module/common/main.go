package module

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"

	_ "github.com/gophab/gophrame/module/common/controller"
	_ "github.com/gophab/gophrame/module/common/service"
)

const (
	MODULE_ID = 102304
)

var _module = &module.Module{
	Name:        "Common",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name)
	module.RegisterModule(_module)
}
