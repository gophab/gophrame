package bootstrap

import (
	// system initialization
	_ "github.com/wjshen/gophrame/controller/api"
	_ "github.com/wjshen/gophrame/core/starter"
	_ "github.com/wjshen/gophrame/security"

	"github.com/wjshen/gophrame/core/logger"

	// init router
	"github.com/wjshen/gophrame/router"
)

// Lazy init
func Init() {
	logger.Info("Initializing Bootstrap...")

	// 1. 初始化 项目根路径，参见 variable 常量包，相关路径：app\global\variable\config.go

	// 初始化服务、路由
	router.InitRouters()

	logger.Info("Initialized Bootstrap")
}
