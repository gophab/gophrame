package wxma

import (
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxma/config"
)

func Start() {
	if config.Setting.Enabled {
		if service, _ := initWxmaService(); service != nil {
			social.RegisterSocialService("wm", service)
		}
	}
}

func initWxmaService() (*WxmaService, error) {
	logger.Info("Initializing WxmaService...")
	return &WxmaService{}, nil
}
