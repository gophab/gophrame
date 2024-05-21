package starter

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/feishu"
	"github.com/wjshen/gophrame/core/social/feishu/config"
)

var feishuService *feishu.FeishuService
var feishuController *feishu.FeishuController

func init() {
	if config.Setting.Enabled {
		if service, _ := initFeishuService(); service != nil {
			social.RegisterSocialService("fs", service)
		}
	}
}

func initFeishuService() (*feishu.FeishuService, error) {
	logger.Info("Initializing FeishuService...")
	feishuService = &feishu.FeishuService{}
	feishuController = &feishu.FeishuController{FeishuService: feishuService}

	inject.InjectValue("feishuController", feishuController)
	return feishuService, nil
}
