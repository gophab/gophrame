package auth

import (
	"github.com/wjshen/gophrame/core/inject"

	AuthModel "github.com/wjshen/gophrame/default/domain/auth"
	AuthRepository "github.com/wjshen/gophrame/default/repository/auth"

	"github.com/wjshen/gophrame/service"
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
