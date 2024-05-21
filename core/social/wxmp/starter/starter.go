package starter

import (
	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxmp"
	"github.com/wjshen/gophrame/core/social/wxmp/config"
)

var wxmpService *wxmp.WxmpService
var wxmpController *wxmp.WxmpController

func init() {
	if config.Setting.Enabled {
		if service, _ := initWxmpService(); service != nil {
			social.RegisterSocialService("mp", service)
		}
	}
}

func initWxmpService() (*wxmp.WxmpService, error) {
	logger.Info("Initializing WxmpService...")
	wxmpService = &wxmp.WxmpService{}
	wxmpController = &wxmp.WxmpController{WxmpService: wxmpService}

	inject.InjectValue("wxmpController", wxmpController)
	return wxmpService, nil
}
