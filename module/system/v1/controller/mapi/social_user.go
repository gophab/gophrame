package mapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/v1/service"
)

type SocialUserMController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService_v1"`
}

var socialUserMController *SocialUserMController = &SocialUserMController{}

func init() {
	inject.InjectValue("socialUserMController_v1", socialUserMController)
}
