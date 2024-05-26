package api

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/default/service"
)

type SocialUserController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserController *SocialUserController = &SocialUserController{}

func init() {
	inject.InjectValue("socialUserController", socialUserController)
}
