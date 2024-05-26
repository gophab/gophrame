package casbin

import (
	"github.com/gophab/gophrame/core/casbin/config"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing Casbin ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		if enforcer, err := InitCasbinEnforcer(); err != nil {
			logger.Error("Load Casbin Enforcer Error: ", err.Error())
		} else if enforcer != nil {
			global.Enforcer = enforcer

			// inject
			logger.Debug("Injected Enforcer")
			inject.InjectValue("enforcer", enforcer)
			logger.Info("Casbin initialized OK")
		}
	}
}
