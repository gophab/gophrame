package auth

import (
	"github.com/gophab/gophrame/core/inject"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthRepository "github.com/gophab/gophrame/default/repository/auth"

	"github.com/gophab/gophrame/service"
)

type RoleMenuService struct {
	service.BaseService
	RoleMenuRepository *AuthRepository.RoleMenuRepository `inject:"roleMenuRepository"`
}

var roleMenuService *RoleMenuService = &RoleMenuService{}

func init() {
	inject.InjectValue("roleMenuService", roleMenuService)
}

func (s *RoleMenuService) GetByRoleIds(ids []string) []AuthModel.RoleMenu {
	return s.RoleMenuRepository.GetByRoleIds(ids)
}
