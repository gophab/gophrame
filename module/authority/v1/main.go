package slink

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"

	// 增加业务API
	_ "github.com/gophab/gophrame/module/authority/v1/controller"
)

var _module = &module.Module{
	Name:        "Authority",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name)
	module.RegisterModule(_module)
}
