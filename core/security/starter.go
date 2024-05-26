package security

import (
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"

	"github.com/gophab/gophrame/core/security/server"
	"github.com/gophab/gophrame/core/security/token"
)

func init() {
	starter.RegisterInitializor(Init)
}

/**
 * 安全框架启动
 */
func Init() {
	logger.Info("Initializing GOES Security Starter")
	token.Init()
	server.Init()
}
