package auth

import (
	"github.com/wjshen/gophrame/core/inject"
	AuthModel "github.com/wjshen/gophrame/domain/auth"

	"github.com/wjshen/gophrame/repository/auth"
	"github.com/wjshen/gophrame/service"
)

type RoleMenuService struct {
	service.BaseService
	RoleMenuRepository *auth.RoleMenuRepository `inject:"roleMenuRepository"`
}

var roleMenuService *RoleMenuService = &RoleMenuService{}

func init() {
	inject.InjectValue("roleMenuService", roleMenuService)
}

func (s *RoleMenuService) GetByRoleIds(ids []string) []AuthModel.RoleMenu {
	return s.RoleMenuRepository.GetByRoleIds(ids)
}
