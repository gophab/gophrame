package openapi

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/service"
)

type SocialUserOpenController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserOpenController *SocialUserOpenController = &SocialUserOpenController{}

func init() {
	inject.InjectValue("socialUserOpenController", socialUserOpenController)
}
