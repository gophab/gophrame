package wxmp

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxmp/config"
)

var wxmpService *WxmpService
var wxmpController *WxmpController

func Start() {
	if config.Setting.Enabled {
		if service, _ := initWxmpService(); service != nil {
			social.RegisterSocialService("mp", service)
		}
	}
}

func initWxmpService() (*WxmpService, error) {
	logger.Info("Initializing WxmpService...")
	wxmpService = &WxmpService{}
	wxmpController = &WxmpController{WxmpService: wxmpService}

	inject.InjectValue("wxmpController", wxmpController)
	return wxmpService, nil
}
