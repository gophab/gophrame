package openapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/v1/service"
)

type SocialUserOpenController struct {
	controller.ResourceController
	SocialUserService *service.SocialUserService `inject:"socialUserService_v1"`
}

var socialUserOpenController *SocialUserOpenController = &SocialUserOpenController{}

func init() {
	inject.InjectValue("socialUserOpenController_v1", socialUserOpenController)
}
