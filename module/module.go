package module

import (
	"github.com/gophab/gophrame/core/logger"

	// init 短链接模块
	_ "github.com/gophab/gophrame/module/slink"

	// init 通用模块
	_ "github.com/gophab/gophrame/module/common"

	// init 权限模块
	_ "github.com/gophab/gophrame/module/operation"

	// init 系统模块
	_ "github.com/gophab/gophrame/module/system"

	// init 授权模块
	_ "github.com/gophab/gophrame/module/authority"
)

func init() {
	logger.Info("Initializing Modules...")
}
