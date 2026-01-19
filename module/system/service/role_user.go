package service

import (
	"errors"

	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/array"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
)

type RoleUserService struct {
	service.BaseService
	RoleUserRepository *repository.RoleUserRepository `inject:"roleUserRepository"`
	RoleRepository     *repository.RoleRepository     `inject:"roleRepository"`
	UserService        *UserService                   `inject:"userService"`
}

var roleUserService *RoleUserService = &RoleUserService{}

func init() {
	inject.InjectValue("roleUserService", roleUserService)

	eventbus.RegisterEventListener("USER_CREATED", roleUserService.onUserCreated)
	eventbus.RegisterEventListener("USER_UPDATED", roleUserService.onUserUpdated)
}

func (s *RoleUserService) ListMembers(roleId, search, tenantId string, pageable query.Pageable) (int64, []*domain.RoleMember) {
	return s.RoleUserRepository.ListMembers(roleId, search, tenantId, pageable)
}

func (s *RoleUserService) ListUsers(roleId, search, tenantId string, pageable query.Pageable) (int64, []*domain.User) {
	return s.RoleUserRepository.ListUsers(roleId, search, tenantId, pageable)
}

// 根据用户id查询所有可能的岗位节点id
func (s *RoleUserService) GetUserRoleIds(userId string) []string {
	//获取用户的所有岗位id
	roleUsers := s.RoleUserRepository.GetByUserId(userId)

	roleIds := []string{}
	for _, v := range roleUsers {
		roleIds = append(roleIds, v.RoleId)
	}

	return roleIds
}

// 根据用户id查询所有可能的岗位节点id
func (s *RoleUserService) GetUserRoles(userId string) []*domain.RoleUser {
	//获取用户的所有岗位id
	roleUsers := s.RoleUserRepository.GetByUserId(userId)
	return roleUsers
}

func (s *RoleUserService) ClearUserRoles(userId string) {
	//获取用户的所有岗位id
	s.RoleUserRepository.DeleteByUserId(userId)
}

func (s *RoleUserService) GetRoleUserIds(roleId string) []string {
	//获取用户的所有岗位id
	roleUsers := s.RoleUserRepository.GetByRoleId(roleId)

	roleIds := []string{}
	for _, v := range roleUsers {
		roleIds = append(roleIds, v.RoleId)
	}

	return roleIds
}

func (s *RoleUserService) AddRoleUserIds(roleId string, userIds []string) ([]*domain.RoleUser, error) {
	//获取用户的所有岗位id
	var roleUsers []*domain.RoleUser = make([]*domain.RoleUser, 0)
	for _, userId := range userIds {
		if user, err := s.UserService.GetById(userId); err == nil {
			roleUser := &domain.RoleUser{
				RoleId: roleId,
				UserId: userId,
			}
			roleUser.TenantId = user.TenantId

			roleUsers = append(roleUsers, roleUser)
		}
	}
	for _, roleUser := range roleUsers {
		if res := s.RoleUserRepository.FirstOrCreate(roleUser); res.Error != nil {
			return nil, res.Error
		}
		eventbus.PublishEvent("USER_ROLE_ADDED", roleUser)
	}
	return roleUsers, nil
}

func (s *RoleUserService) DeleteRoleUserIds(roleId string, userIds []string) error {
	//获取用户的所有岗位id
	for _, userId := range userIds {
		if b := s.RoleUserRepository.DeleteData(roleId, userId); !b {
			return errors.New("Delete Error")
		}
		eventbus.PublishEvent("USER_ROLE_DELETED", &domain.RoleUser{
			RoleId: roleId,
			UserId: userId,
		})
	}
	return nil
}

func (s *RoleUserService) AddUserRoleIds(userId string, roleIds []string) ([]*domain.RoleUser, error) {
	//获取用户的所有岗位id
	if user, err := s.UserService.GetById(userId); err == nil {
		var roleUsers []*domain.RoleUser = make([]*domain.RoleUser, 0)
		for _, roleId := range roleIds {
			roleUser := &domain.RoleUser{
				RoleId: roleId,
				UserId: userId,
			}
			roleUser.TenantId = user.TenantId

			roleUsers = append(roleUsers, roleUser)
		}
		for _, roleUser := range roleUsers {
			if res := s.RoleUserRepository.FirstOrCreate(roleUser); res.Error != nil {
				return nil, res.Error
			}
			eventbus.PublishEvent("USER_ROLE_ADDED", roleUser)
		}
		return roleUsers, nil
	}
	return nil, nil
}

func (s *RoleUserService) AddUserRoles(userId string, roles []string) ([]*domain.RoleUser, error) {
	//获取用户的所有岗位id
	rs, err := s.RoleRepository.GetRoles(core.M{"names": roles})
	if err != nil || len(rs) == 0 {
		return []*domain.RoleUser{}, err
	}

	roleIds := array.Map(rs, func(r *domain.Role) string {
		return r.Id
	})

	return s.AddUserRoleIds(userId, roleIds)
}

func (s *RoleUserService) DeleteUserRoleIds(userId string, roleIds []string) error {
	//获取用户的所有岗位id
	for _, roleId := range roleIds {
		if b := s.RoleUserRepository.DeleteData(roleId, userId); !b {
			return errors.New("Delete Error")
		}
		eventbus.PublishEvent("USER_ROLE_DELETED", &domain.RoleUser{
			RoleId: roleId,
			UserId: userId,
		})
	}
	return nil
}

func (s *RoleUserService) DeleteUserRoles(userId string, roles []string) error {
	//获取用户的所有岗位id
	rs, err := s.RoleRepository.GetRoles(core.M{"names": roles})
	if err != nil || len(rs) == 0 {
		return err
	}

	roleIds := array.Map(rs, func(r *domain.Role) string {
		return r.Id
	})

	return s.DeleteUserRoleIds(userId, roleIds)
}

func (s *RoleUserService) HasRole(userId string, role string) bool {
	if results := s.RoleUserRepository.GetUserRoles(userId); len(results) > 0 {
		for _, result := range results {
			if result.Id == role || result.Name == role {
				return true
			}
		}
	}
	return false
}

func (s *RoleUserService) HasRoleId(userId string, roleId string) bool {
	return s.RoleUserRepository.GetByUserIdAndRoleId(userId, roleId) != nil
}

func (s *RoleUserService) onUserCreated(event string, args ...any) {
	var user = args[0].(*domain.User)

	switch user.TenantId {
	case "SYSTEM":
		// 1. Default add to ROLE:00000000000001 - 系统用户
		role, err := s.RoleRepository.GetByName("operator", "SYSTEM")
		if err == nil && role != nil {
			// 自动添加到系统用户角色
			s.AddRoleUserIds(role.Id, []string{user.Id})
		}
	case "PUBLIC":
		// 1. Default add to ROLE:00000000000001 - 系统用户
		role, err := s.RoleRepository.GetByName("user", "SYSTEM")
		if err == nil && role != nil {
			// 自动添加到系统用户角色
			s.AddRoleUserIds(role.Id, []string{user.Id})
		}
	default:
		// 2. Default add to ROLE:00000000000003 - 企业用户
		role, err := s.RoleRepository.GetByName("member", "SYSTEM")
		if err == nil && role != nil {
			// 自动添加到企业用户角色
			s.AddRoleUserIds(role.Id, []string{user.Id})
		}
	}

	// 创建时设置角色
	if len(user.Roles) > 0 {
		var roleIds = make([]string, 0)
		for _, role := range user.Roles {
			if len(role.Id) > 0 {
				roleIds = append(roleIds, role.Id)
			} else if len(role.Name) > 0 {
				if r, err := s.RoleRepository.GetByName(role.Name, user.TenantId); err == nil && r != nil {
					roleIds = append(roleIds, r.Id)
				}
			}
		}

		if len(roleIds) > 0 {
			s.AddUserRoleIds(user.Id, roleIds)
		}
	}
}

func (s *RoleUserService) onUserUpdated(event string, args ...any) {
	var user = args[0].(*domain.User)

	// 创建时设置角色
	if user.Roles != nil {
		// 清除所有
		s.ClearUserRoles(user.Id)

		if len(user.Roles) > 0 {
			var roleIds = make([]string, 0)
			for _, role := range user.Roles {
				if len(role.Id) > 0 {
					roleIds = append(roleIds, role.Id)
				} else if len(role.Name) > 0 {
					if r, err := s.RoleRepository.GetByName(role.Name, user.TenantId); err == nil && r != nil {
						roleIds = append(roleIds, r.Id)
					}
				}
			}

			s.AddUserRoleIds(user.Id, roleIds)
		}
	}
}
