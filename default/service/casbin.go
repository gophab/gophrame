package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/default/domain"
)

// CasbinService负责更新Casbin Enforce数据至
type CasbinService struct {
	Enforcer    *casbin.SyncedEnforcer `inject:"enforcer"`
	UserService *UserService           `inject:"userService"`
	RoleService *RoleService           `inject:"roleService"`
}

var casbinService = &CasbinService{}

func init() {
	inject.InjectValue("casbinService", casbinService)

	eventbus.RegisterEventListener("ROLE_CREATED", casbinService.onRoleEvent)
	eventbus.RegisterEventListener("ROLE_UPDATED", casbinService.onRoleEvent)
	eventbus.RegisterEventListener("ROLE_DELETED", casbinService.onRoleEvent)

	eventbus.RegisterEventListener("USER_CREATED", casbinService.onUserEvent)
	eventbus.RegisterEventListener("USER_DELETED", casbinService.onUserEvent)
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
		role, err := s.RoleService.Get(id)
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
	case "USER_CREATED":
		s.refreshUserPolicy(user)
		break
	}
}

func (s *CasbinService) onRoleEvent(event string, args ...interface{}) {
	role := args[0].(*domain.Role)
	switch event {
	case "ROLE_CREATED":
		break
	case "ROLE_UPDATED":
		break
	case "ROLE_DELETED":
		s.Enforcer.DeletePermissionsForUser(role.Name)
		break
	}
}
