package service

import (
	"errors"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
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

func (s *RoleService) CreateRole(role *domain.Role) (*domain.Role, error) {
	name, _ := s.RoleRepository.CheckRoleName(role.Name, role.TenantId)
	if name {
		return nil, errors.New("name 名字重复,请更改！")
	}

	res, err := s.RoleRepository.CreateRole(role)
	if err != nil {
		return nil, err
	}

	eventbus.PublishEvent("ROLE_CREATED", res)

	return res, nil
}

func (s *RoleService) UpdateRole(role *domain.Role) (*domain.Role, error) {
	name, _ := s.RoleRepository.CheckRoleNameId(role.Name, role.Id, role.TenantId)
	if name {
		return nil, errors.New("name 名字重复,请更改！")
	}

	role, err := s.RoleRepository.UpdateRole(role)
	if err != nil {
		return nil, err
	}

	eventbus.PublishEvent("ROLE_UPDATED", role)
	return role, nil
}

func (s *RoleService) PatchRole(role *domain.Role, data map[string]interface{}) (*domain.Role, error) {
	if n, b := data["name"]; b {
		name, _ := s.RoleRepository.CheckRoleNameId(n.(string), role.Id, role.TenantId)
		if name {
			return nil, errors.New("name 名字重复,请更改！")
		}
	}

	role, err := s.RoleRepository.PatchRole(role.Id, data)
	if err != nil {
		return nil, err
	}

	eventbus.PublishEvent("ROLE_UPDATED", role)
	return role, nil
}

func (s *RoleService) DeleteById(id string) error {
	role, err := s.RoleRepository.GetById(id)
	if err != nil {
		return err
	}

	if role != nil {
		err := s.RoleRepository.DeleteById(id)
		if err != nil {
			return err
		}

		eventbus.PublishEvent("ROLE_DELETED", role)
	}

	return nil
}

func (s *RoleService) GetById(id string) (*domain.Role, error) {
	role, err := s.RoleRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetUserRoles(userId string) ([]*domain.Role, error) {
	var conds = make(map[string]interface{})
	conds["user_id"] = userId
	return s.RoleRepository.GetRoles(conds)
}

func (s *RoleService) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Role, error) {
	return s.RoleRepository.FindRoles(conds, pageable)
}

func (s *RoleService) FindAvailable(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Role, error) {
	return s.RoleRepository.FindAvailableRoles(conds, pageable)
}

func (s *RoleService) ExistById(id string) (bool, error) {
	return s.RoleRepository.ExistById(id)
}

func (s *RoleService) Count(conds map[string]interface{}) (int64, error) {
	return s.RoleRepository.GetRoleTotal(conds)
}
