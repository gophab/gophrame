package wxma

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social"
	SocialConfig "github.com/gophab/gophrame/core/social/config"
	"github.com/gophab/gophrame/core/social/wxma/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Info("Initializing WxmaService: ...", SocialConfig.Setting.Enabled && config.Setting.Enabled)
	if SocialConfig.Setting.Enabled && config.Setting.Enabled {
		wxmaService := &WxmaService{}
		inject.InjectValue("wxmaService", wxmaService)

		social.RegisterSocialService("wma", wxmaService)
	}
}
