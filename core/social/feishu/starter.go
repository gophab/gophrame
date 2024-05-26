package feishu

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social"
	SocialConfig "github.com/gophab/gophrame/core/social/config"
	"github.com/gophab/gophrame/core/social/feishu/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Initializing FeishuService: ...", SocialConfig.Setting.Enabled && config.Setting.Enabled)
	if SocialConfig.Setting.Enabled && config.Setting.Enabled {
		feishuService := &FeishuService{}
		inject.InjectValue("feishuService", feishuService)

		social.RegisterSocialService("fs", feishuService)

		feishuController := &FeishuController{FeishuService: feishuService}

		inject.InjectValue("feishuController", feishuController)

		controller.AddController(feishuController)
	}
}
