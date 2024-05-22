package wxcp

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxcp/config"
)

var wxcpService *WxcpService
var wxcpController *WxcpController

func Start() {
	if config.Setting.Enabled {
		if service, _ := initWxcpService(); service != nil {
			social.RegisterSocialService("ww", service)
		}
	}
}

func initWxcpService() (*WxcpService, error) {
	logger.Info("Initializing WxcpService...")
	wxcpService = &WxcpService{}
	wxcpController = &WxcpController{WxcpService: wxcpService}

	inject.InjectValue("wxcpController", wxcpController)
	return wxcpService, nil
}
