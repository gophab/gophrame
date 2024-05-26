package wxcp

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social"
	SocialConfig "github.com/gophab/gophrame/core/social/config"
	"github.com/gophab/gophrame/core/social/wxcp/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing WxcpService: ...", SocialConfig.Setting.Enabled && config.Setting.Enabled)
	if SocialConfig.Setting.Enabled && config.Setting.Enabled {
		wxcpService := &WxcpService{}

		social.RegisterSocialService("ww", wxcpService)

		wxcpController := &WxcpController{WxcpService: wxcpService}
		inject.InjectValue("wxcpController", wxcpController)

		controller.AddController(wxcpController)
	}
}
