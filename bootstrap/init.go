package bootstrap

import (
	// system initialization
	_ "github.com/gophab/gophrame/core/destroy" // 监听程序退出信号，用于资源的释放
	_ "github.com/gophab/gophrame/core/engine"
	_ "github.com/gophab/gophrame/core/eventbus"
	_ "github.com/gophab/gophrame/core/security"
	_ "github.com/gophab/gophrame/core/snowflake"

	// core
	"github.com/gophab/gophrame/core/command"

	_ "github.com/gophab/gophrame/core/casbin"
	_ "github.com/gophab/gophrame/core/email"
	_ "github.com/gophab/gophrame/core/email/code"
	_ "github.com/gophab/gophrame/core/i18n"
	_ "github.com/gophab/gophrame/core/microservice"
	_ "github.com/gophab/gophrame/core/mongo"

	_ "github.com/gophab/gophrame/core/rabbitmq"
	_ "github.com/gophab/gophrame/core/redis"
	_ "github.com/gophab/gophrame/core/sensitive"
	_ "github.com/gophab/gophrame/core/server"
	_ "github.com/gophab/gophrame/core/sms"
	_ "github.com/gophab/gophrame/core/sms/code"
	_ "github.com/gophab/gophrame/core/websocket"

	// starter
	_ "github.com/gophab/gophrame/core/database/starter"
	_ "github.com/gophab/gophrame/core/identify/starter"
	_ "github.com/gophab/gophrame/core/oss/starter"
	_ "github.com/gophab/gophrame/core/payment/starter"
	_ "github.com/gophab/gophrame/core/social/starter"

	// system core
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/router"
	"github.com/gophab/gophrame/core/starter"

	// init config
	"github.com/gophab/gophrame/config"

	// init service
	_ "github.com/gophab/gophrame/service"

	// init controller
	_ "github.com/gophab/gophrame/controller"

	// init moudles
	_ "github.com/gophab/gophrame/module"
)

// Lazy init
func Init() {
	// 0. init()

	// 1. Register() - RegisterConfig() - RegisterInitializor() - RegisterStarter() - RegisterTerminater - RegisterPlugin - RegisterPlugin
	logger.Info("Initializing Framework Bootstrap...")

	// 1. 解析命令行参数
	command.Init()

	// 2. 读取配置
	config.Init()

	// 3. 启动器
	starter.Init()

	// 4. 启动router
	router.Init()

	// 5. 启动器
	starter.Start()

	logger.Info("Initialized Framework Bootstrap")
}
