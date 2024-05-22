package router

import (
	DefaultRouter "github.com/wjshen/gophrame/core/router"

	"github.com/wjshen/gophrame/controller"
)

func Init() {
	// 初始化缺省路由
	root := DefaultRouter.InitDefaultRouters()

	// API
	controller.InitRouter(root)
}
