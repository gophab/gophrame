package websocket

import (
	"github.com/wjshen/gophrame/core/websocket/config"

	"github.com/wjshen/gophrame/core/global"
	"github.com/wjshen/gophrame/core/logger"
)

func init() {
	// websocket Hub中心启动
	if config.Setting.Enabled {
		logger.Info("Enable websocket")
		// websocket 管理中心hub全局初始化一份
		global.WebsocketHub = CreateHubFactory()
		if WebsocketHub, ok := global.WebsocketHub.(*Hub); ok {
			go WebsocketHub.Run()
			logger.Info("Running websocket hub OK")
		}
	}
}
