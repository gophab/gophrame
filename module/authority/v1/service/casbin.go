package service

import (
	"github.com/casbin/casbin/v2"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/v1/domain"
	"github.com/gophab/gophrame/module/system/v1/service"
)

// CasbinService负责更新Casbin Enforce数据至
type CasbinService struct {
	Enforcer    *casbin.SyncedEnforcer `inject:"enforcer_v1"`
	UserService *service.UserService   `inject:"userService_v1"`
	RoleService *service.RoleService   `inject:"roleService_v1"`
}

var casbinService = &CasbinService{}

func init() {
	inject.InjectValue("casbinService_v1", casbinService)

	eventbus.RegisterEventListener("ROLE_CREATED_v1", casbinService.onRoleEvent)
	eventbus.RegisterEventListener("ROLE_UPDATED_v1", casbinService.onRoleEvent)
	eventbus.RegisterEventListener("ROLE_DELETED_v1", casbinService.onRoleEvent)

	eventbus.RegisterEventListener("USER_CREATED_v1", casbinService.onUserEvent)
	eventbus.RegisterEventListener("USER_DELETED_v1", casbinService.onUserEvent)
}

// LoadAllPolicy 加载所有的角色策略
func (s *CasbinService) LoadAllPolicy() error {
	if s.Enforcer != nil {
	}

	return nil
}

// LoadPolicy 加载角色权限策略
func (s *CasbinService) LoadPolicy(id string) error {
	if s.Enforcer != nil {
		role, err := s.RoleService.GetById(id)
		if err != nil {
			return err
		}

		return s.loadPolicy(role)
	}

	return nil
}

// LoadPolicy 加载角色权限策略
func (s *CasbinService) loadPolicy(role *domain.Role) error {
	if s.Enforcer != nil {
		s.Enforcer.DeleteRole(role.Name)
	}

	return nil
}

func (s *CasbinService) refreshUserPolicy(user *domain.User) {
	// 转换User信息到Casbin表
	if _, err := s.RoleService.GetUserRoles(user.Id); err == nil {

		if s.Enforcer != nil {
			s.Enforcer.LoadPolicy()
		}
	}

}

func (s *CasbinService) onUserEvent(event string, args ...interface{}) {
	var user = args[0].(*domain.User)

	switch event {
	case "USER_CREATED_v1":
		s.refreshUserPolicy(user)
	}
}

func (s *CasbinService) onRoleEvent(event string, args ...interface{}) {
	role := args[0].(*domain.Role)
	switch event {
	case "ROLE_CREATED_v1":
		break
	case "ROLE_UPDATED_v1":
		break
	case "ROLE_DELETED_v1":
		if s.Enforcer != nil {
			s.Enforcer.DeletePermissionsForUser(role.Name)
		}
	}
}
