package service

import (
	"errors"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/query"

	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/repository"
	"github.com/wjshen/gophrame/service/dto"

	"github.com/casbin/casbin/v2"
)

type RoleService struct {
	BaseService
	RoleResposity *repository.RoleRepository `inject:"roleRepository"`
	Enforcer      *casbin.SyncedEnforcer     `inject:"enforcer"`
}

var roleService *RoleService = &RoleService{}

func GetRoleService() *RoleService {
	return roleService
}

func init() {
	logger.Debug("Initializing RoleService...")
	logger.Debug("Inject roleService")
	inject.InjectValue("roleService", roleService)
	logger.Debug("Initialized RoleService")
}

func (s *RoleService) Add(role *dto.Role) (*domain.Role, error) {
	name, _ := s.RoleResposity.CheckRoleName(role.Name)
	if name {
		return nil, errors.New("name 名字重复,请更改！")
	}

	res, err := s.RoleResposity.AddRole(map[string]interface{}{
		"name": role.Name,
	})

	if err != nil {
		return nil, err
	}

	err = s.LoadPolicy(role.Id)
	if err != nil {
		return res, errors.New("load policy failed")
	}

	return res, nil
}

func (s *RoleService) Edit(role *dto.Role) error {
	name, _ := s.RoleResposity.CheckRoleNameId(role.Name, role.Id)
	if name {
		return errors.New("name 名字重复,请更改！")
	}

	err := s.RoleResposity.EditRole(role.Id, map[string]interface{}{
		"name": role.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *RoleService) Get(id string) (*domain.Role, error) {
	role, err := s.RoleResposity.GetRole(id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetAll(role *dto.Role, pageable query.Pageable) ([]*domain.Role, error) {
	if role.Id != "" {
		maps := make(map[string]interface{})
		maps["del_flag"] = false
		maps["id"] = role.Id

		roles, err := s.RoleResposity.GetRoles(maps, pageable)
		if err != nil {
			return nil, err
		}

		return roles, nil
	} else {
		roles, err := s.RoleResposity.GetRoles(role.GetMaps(), pageable)
		if err != nil {
			return nil, err
		}
		return roles, nil
	}
}

func (s *RoleService) Delete(id string) error {
	role, err := s.RoleResposity.GetRole(id)
	if err != nil {
		return err
	}

	if role != nil {
		err := s.RoleResposity.DeleteRole(id)
		if err != nil {
			return err
		}
		s.Enforcer.DeletePermissionsForUser(role.Name)
	}

	return nil
}

func (s *RoleService) ExistByID(id string) (bool, error) {
	return s.RoleResposity.ExistRoleByID(id)
}

func (s *RoleService) Count(role *dto.Role) (int64, error) {
	return s.RoleResposity.GetRoleTotal(role.GetMaps())
}

// LoadAllPolicy 加载所有的角色策略
func (s *RoleService) LoadAllPolicy() error {
	if s.Enforcer != nil {
		roles, err := s.RoleResposity.GetRolesAll()
		if err != nil {
			return err
		}

		for _, role := range roles {
			err = s.loadPolicy(role)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadPolicy 加载角色权限策略
func (s *RoleService) LoadPolicy(id string) error {
	if s.Enforcer != nil {
		role, err := s.RoleResposity.GetRole(id)
		if err != nil {
			return err
		}

		return s.loadPolicy(role)
	}

	return nil
}

// LoadPolicy 加载角色权限策略
func (s *RoleService) loadPolicy(role *domain.Role) error {
	if s.Enforcer != nil {
		s.Enforcer.DeleteRole(role.Name)
	}

	return nil
}
