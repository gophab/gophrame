package api

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"

	"github.com/wjshen/gophrame/default/service"
)

type SocialUserController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserController *SocialUserController = &SocialUserController{}

func init() {
	inject.InjectValue("socialUserController", socialUserController)
}
