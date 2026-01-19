package wxmp

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social"
	SocialConfig "github.com/gophab/gophrame/core/social/config"
	"github.com/gophab/gophrame/core/social/wxmp/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing WxmpService: ...", SocialConfig.Setting.Enabled && config.Setting.Enabled)
	if SocialConfig.Setting.Enabled && config.Setting.Enabled {
		wxmpService := &WxmpService{}
		inject.InjectValue("wxmpService", wxmpService)

		social.RegisterSocialService("wx", wxmpService)

		wxmpController := &WxmpController{WxmpService: wxmpService}
		inject.InjectValue("wxmpController", wxmpController)

		controller.AddController(wxmpController)
	}
}
