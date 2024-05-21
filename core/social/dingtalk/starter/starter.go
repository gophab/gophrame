package starter

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/dingtalk"
	"github.com/wjshen/gophrame/core/social/dingtalk/config"
)

var dingtalkService *dingtalk.DingtalkService
var dingtalkController *dingtalk.DingtalkController

func init() {
	if config.Setting.Enabled {
		if service, _ := initDingtalkService(); service != nil {
			social.RegisterSocialService("dt", service)
		}
	}
}

func initDingtalkService() (*dingtalk.DingtalkService, error) {
	logger.Info("Initializing DingtalkService...")
	dingtalkService = &dingtalk.DingtalkService{}
	dingtalkController = &dingtalk.DingtalkController{DingtalkService: dingtalkService}

	inject.InjectValue("dingtalkController", dingtalkController)
	return dingtalkService, nil
}
