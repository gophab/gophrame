package bootstrap

import (
	// system initialization
	_ "github.com/wjshen/gophrame/core/casbin"
	_ "github.com/wjshen/gophrame/core/database"
	_ "github.com/wjshen/gophrame/core/destroy" // 监听程序退出信号，用于资源的释放
	_ "github.com/wjshen/gophrame/core/email"
	_ "github.com/wjshen/gophrame/core/email/code"
	_ "github.com/wjshen/gophrame/core/engine"
	_ "github.com/wjshen/gophrame/core/eventbus"
	_ "github.com/wjshen/gophrame/core/microservice"
	_ "github.com/wjshen/gophrame/core/rabbitmq"
	_ "github.com/wjshen/gophrame/core/redis"
	_ "github.com/wjshen/gophrame/core/security"
	_ "github.com/wjshen/gophrame/core/sms"
	_ "github.com/wjshen/gophrame/core/sms/code"
	_ "github.com/wjshen/gophrame/core/snowflake"
	_ "github.com/wjshen/gophrame/core/social/starter"
	_ "github.com/wjshen/gophrame/core/websocket"

	// _ "github.com/wjshen/gophrame/security"

	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/starter"

	// init config
	"github.com/wjshen/gophrame/config"

	// init router
	"github.com/wjshen/gophrame/router"
)

// Lazy init
func Init() {
	logger.Info("Initializing Bootstrap...")

	// 1. 读取配置
	config.Init()

	// 2. 启动器
	starter.Start()

	// 3. 初始化服务、路由
	router.Init()

	logger.Info("Initialized Bootstrap")
}
