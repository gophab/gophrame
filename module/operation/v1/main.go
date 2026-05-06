package slink

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"

	// 增加业务API
	_ "github.com/gophab/gophrame/module/operation/v1/controller"
)

const (
	MODULE_ID = 10110
)

var _module = &module.Module{
	Name:        "OperationV1",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name, "v1")
	module.RegisterModule(_module)
}
