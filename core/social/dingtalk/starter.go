package dingtalk

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/dingtalk/config"
)

var dingtalkService *DingtalkService
var dingtalkController *DingtalkController

func Start() {
	if config.Setting.Enabled {
		if service, _ := initDingtalkService(); service != nil {
			social.RegisterSocialService("dt", service)
		}
	}
}

func initDingtalkService() (*DingtalkService, error) {
	logger.Info("Initializing DingtalkService...")
	dingtalkService = &DingtalkService{}
	dingtalkController = &DingtalkController{DingtalkService: dingtalkService}

	inject.InjectValue("dingtalkController", dingtalkController)
	return dingtalkService, nil
}
