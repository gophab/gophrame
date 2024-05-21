package security

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/logger"

	_ "github.com/wjshen/gophrame/core/security/server"
	_ "github.com/wjshen/gophrame/core/security/token"
)

/**
 * 安全框架启动
 */
func init() {
	logger.Info("Initializing GOES Security Starter")
}
