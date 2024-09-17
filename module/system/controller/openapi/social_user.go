package openapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/service"
)

type SocialUserOpenController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
}

var socialUserOpenController *SocialUserOpenController = &SocialUserOpenController{}

func init() {
	inject.InjectValue("socialUserOpenController", socialUserOpenController)
}
