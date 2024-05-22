package feishu

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/feishu/config"
)

var feishuService *FeishuService
var feishuController *FeishuController

func Start() {
	if config.Setting.Enabled {
		if service, _ := initFeishuService(); service != nil {
			social.RegisterSocialService("fs", service)
		}
	}
}

func initFeishuService() (*FeishuService, error) {
	logger.Info("Initializing FeishuService...")
	feishuService = &FeishuService{}
	feishuController = &FeishuController{FeishuService: feishuService}

	inject.InjectValue("feishuController", feishuController)
	return feishuService, nil
}
