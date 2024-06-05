package service

import (
	"errors"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
	"github.com/gophab/gophrame/default/service/dto"
)

type RoleService struct {
	service.BaseService
	RoleRepository *repository.RoleRepository `inject:"roleRepository"`
}

var roleService *RoleService = &RoleService{}

func init() {
	logger.Debug("Initializing RoleService...")
	logger.Debug("Inject roleService")
	inject.InjectValue("roleService", roleService)
	logger.Debug("Initialized RoleService")
}

func (s *RoleService) Add(role *dto.Role) (*domain.Role, error) {
	name, _ := s.RoleRepository.CheckRoleName(role.Name)
	if name {
		return nil, errors.New("name 名字重复,请更改！")
	}

	res, err := s.RoleRepository.AddRole(map[string]interface{}{
		"name": role.Name,
	})

	if err != nil {
		return nil, err
	}

	eventbus.PublishEvent("ROLE_CREATED", res)

	return res, nil
}

func (s *RoleService) Edit(role *dto.Role) error {
	name, _ := s.RoleRepository.CheckRoleNameId(role.Name, role.Id)
	if name {
		return errors.New("name 名字重复,请更改！")
	}

	err := s.RoleRepository.EditRole(role.Id, map[string]interface{}{
		"name": role.Name,
	})
	if err != nil {
		return err
	}

	eventbus.PublishEvent("ROLE_UPDATED", role)

	return nil
}

func (s *RoleService) Delete(id string) error {
	role, err := s.RoleRepository.GetRole(id)
	if err != nil {
		return err
	}

	if role != nil {
		err := s.RoleRepository.DeleteRole(id)
		if err != nil {
			return err
		}

		eventbus.PublishEvent("ROLE_DELETED", role)
	}

	return nil
}

func (s *RoleService) Get(id string) (*domain.Role, error) {
	role, err := s.RoleRepository.GetRole(id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetUserRoles(userId string) ([]*domain.Role, error) {
	var params = make(map[string]interface{})
	params["user_id"] = userId
	return s.RoleRepository.FindRolesAll(params)
}

func (s *RoleService) GetAll(role *dto.Role, pageable query.Pageable) ([]*domain.Role, error) {
	if role.Id != "" {
		maps := make(map[string]interface{})
		maps["del_flag"] = false
		maps["id"] = role.Id

		roles, err := s.RoleRepository.FindRoles(maps, pageable)
		if err != nil {
			return nil, err
		}

		return roles, nil
	} else {
		roles, err := s.RoleRepository.FindRoles(role.GetMaps(), pageable)
		if err != nil {
			return nil, err
		}
		return roles, nil
	}
}

func (s *RoleService) ExistByID(id string) (bool, error) {
	return s.RoleRepository.ExistRoleByID(id)
}

func (s *RoleService) Count(role *dto.Role) (int64, error) {
	return s.RoleRepository.GetRoleTotal(role.GetMaps())
}
