package auth

import (
	"github.com/wjshen/gophrame/core/inject"

	AuthModel "github.com/wjshen/gophrame/default/domain/auth"
	AuthRepository "github.com/wjshen/gophrame/default/repository/auth"

	"github.com/wjshen/gophrame/service"
)

type MenuButtonService struct {
	service.BaseService
	MenuButtonRepository *AuthRepository.MenuButtonRepository `inject:"menuButtonRepository"`
}

var menuButtonService *MenuButtonService = &MenuButtonService{}

func init() {
	inject.InjectValue("menuButtonService", menuButtonService)
}

func (s *MenuButtonService) List(fid int64) (int64, []AuthModel.MenuButton) {
	return s.MenuButtonRepository.List(fid)
}

func (s *MenuButtonService) GetByButtonId(buttonId int64) bool {
	return s.MenuButtonRepository.GetByButtonId(buttonId)
}
