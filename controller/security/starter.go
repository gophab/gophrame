package security

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Enable Register: ...", global.EnableRegister)
	if global.EnableRegister {
		controller.AddController(&SecurityController{})
	}
}
