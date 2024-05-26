package websocket

import (
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/websocket/config"

	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Enable websocket: ...", config.Setting.Enabled)
	// websocket Hub中心启动
	if config.Setting.Enabled {
		// websocket 管理中心hub全局初始化一份
		global.WebsocketHub = CreateHubFactory()
		if WebsocketHub, ok := global.WebsocketHub.(*Hub); ok {
			go WebsocketHub.Run()
			logger.Info("Running websocket hub OK")
		}
	}
}
