package slink

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/module"

	// 增加业务配置
	_ "github.com/gophab/gophrame/module/slink/config"

	// 增加业务API
	_ "github.com/gophab/gophrame/module/slink/controller"
)

const (
	MODULE_ID = 102101
)

var _module = &module.Module{
	Name:        "ShortLink",
	Description: "",
}

func init() {
	logger.Info("Register module: ", _module.Name)
	module.RegisterModule(_module)
}
