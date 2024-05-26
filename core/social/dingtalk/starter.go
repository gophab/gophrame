package dingtalk

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social"
	SocialConfig "github.com/gophab/gophrame/core/social/config"
	"github.com/gophab/gophrame/core/social/dingtalk/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing DingtalkService: ...", SocialConfig.Setting.Enabled && config.Setting.Enabled)
	if SocialConfig.Setting.Enabled && config.Setting.Enabled {
		dingtalkService := &DingtalkService{}
		inject.InjectValue("dingtalkService", dingtalkService)

		social.RegisterSocialService("dt", dingtalkService)

		dingtalkController := &DingtalkController{DingtalkService: dingtalkService}
		inject.InjectValue("dingtalkController", dingtalkController)

		controller.AddController(dingtalkController)
	}
}
