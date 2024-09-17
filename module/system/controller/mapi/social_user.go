package mapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/service"
)

type SocialUserMController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserMController *SocialUserMController = &SocialUserMController{}

func init() {
	inject.InjectValue("socialUserMController", socialUserMController)
}
