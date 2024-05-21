package mapi

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/service"
)

type SocialUserMController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserMController *SocialUserMController = &SocialUserMController{}

func init() {
	inject.InjectValue("socialUserMController", socialUserMController)
}
