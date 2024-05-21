package starter

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxma"
	"github.com/wjshen/gophrame/core/social/wxma/config"
)

func init() {
	if config.Setting.Enabled {
		if service, _ := initWxmaService(); service != nil {
			social.RegisterSocialService("wm", service)
		}
	}
}

func initWxmaService() (*wxma.WxmaService, error) {
	logger.Info("Initializing WxmaService...")
	return &wxma.WxmaService{}, nil
}
