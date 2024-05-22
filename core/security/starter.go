package security

import (
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/starter"

	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/security/token"
)

func init() {
	starter.RegisterStarter(Start)
}

/**
 * 安全框架启动
 */
func Start() {
	logger.Info("Initializing GOES Security Starter")
	server.Start()
	token.Start()
}
