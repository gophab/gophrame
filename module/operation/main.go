package slink

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"

	// 增加业务API
	_ "github.com/gophab/gophrame/module/operation/controller"
)

const (
	MODULE_ID = 10110
)

var _module = &module.Module{
	Name:        "Operation",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name)
	module.RegisterModule(_module)
}
