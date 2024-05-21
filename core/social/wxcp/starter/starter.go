package starter

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxcp"
	"github.com/wjshen/gophrame/core/social/wxcp/config"
)

var wxcpService *wxcp.WxcpService
var wxcpController *wxcp.WxcpController

func init() {
	if config.Setting.Enabled {
		if service, _ := initWxcpService(); service != nil {
			social.RegisterSocialService("ww", service)
		}
	}
}

func initWxcpService() (*wxcp.WxcpService, error) {
	logger.Info("Initializing WxcpService...")
	wxcpService = &wxcp.WxcpService{}
	wxcpController = &wxcp.WxcpController{WxcpService: wxcpService}

	inject.InjectValue("wxcpController", wxcpController)
	return wxcpService, nil
}
