package service

import (
	"errors"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
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
}

func (u *RoleUserService) ListMembers(roleId, search, tenantId string, pageable query.Pageable) (int64, []*domain.RoleMember) {
	return u.RoleUserRepository.ListMembers(roleId, search, tenantId, pageable)
}

func (u *RoleUserService) ListUsers(roleId, search, tenantId string, pageable query.Pageable) (int64, []*domain.User) {
	return u.RoleUserRepository.ListUsers(roleId, search, tenantId, pageable)
}

// 根据用户id查询所有可能的岗位节点id
func (u *RoleUserService) GetUserRoleIds(userId string) []string {
	//获取用户的所有岗位id
	roleUsers := u.RoleUserRepository.GetByUserId(userId)

	roleIds := []string{}
	for _, v := range roleUsers {
		roleIds = append(roleIds, v.RoleId)
	}

	return roleIds
}

func (u *RoleUserService) GetRoleUserIds(roleId string) []string {
	//获取用户的所有岗位id
	roleUsers := u.RoleUserRepository.GetByRoleId(roleId)

	roleIds := []string{}
	for _, v := range roleUsers {
		roleIds = append(roleIds, v.RoleId)
	}

	return roleIds
}

func (u *RoleUserService) AddRoleUserIds(roleId string, userIds []string) ([]*domain.RoleUser, error) {
	//获取用户的所有岗位id
	var roleUsers []*domain.RoleUser = make([]*domain.RoleUser, 0)
	for _, userId := range userIds {
		if user, err := u.UserService.GetById(userId); err == nil {
			roleUser := &domain.RoleUser{
				RoleId: roleId,
				UserId: userId,
			}
			roleUser.TenantId = user.TenantId

			roleUsers = append(roleUsers, roleUser)
		}
	}
	for _, roleUser := range roleUsers {
		if res := u.RoleUserRepository.FirstOrCreate(roleUser); res.Error != nil {
			return nil, res.Error
		}
	}
	return roleUsers, nil
}

func (u *RoleUserService) DeleteRoleUserIds(roleId string, userIds []string) error {
	//获取用户的所有岗位id
	for _, userId := range userIds {
		if b := u.RoleUserRepository.DeleteData(roleId, userId); !b {
			return errors.New("Delete Error")
		}
	}
	return nil
}

func (s *RoleUserService) onUserCreated(event string, args ...interface{}) {
	var user = args[0].(*domain.User)

	if user.TenantId == "SYSTEM" {
		// 1. Default add to ROLE:00000000000001 - 系统用户
		role, err := s.RoleRepository.GetByName("operator", "SYSTEM")
		if err == nil && role != nil {
			// 自动添加到系统用户角色
			s.AddRoleUserIds(role.Id, []string{user.Id})
		}
	} else {
		// 2. Default add to ROLE:00000000000003 - 企业用户
		role, err := s.RoleRepository.GetByName("member", "SYSTEM")
		if err == nil && role != nil {
			// 自动添加到企业用户角色
			s.AddRoleUserIds(role.Id, []string{user.Id})
		}
	}
}
