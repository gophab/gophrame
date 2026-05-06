package api

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/v1/service"
)

type SocialUserController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService_v1"`
}

var socialUserController *SocialUserController = &SocialUserController{}

func init() {
	inject.InjectValue("socialUserController_v1", socialUserController)
}
