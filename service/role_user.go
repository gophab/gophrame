package service

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/repository"
)

type RoleUserService struct {
	BaseService
	RoleUserRepository *repository.RoleUserRepository `inject:"roleUserRepository"`
	RoleRepository     *repository.RoleRepository     `inject:"roleRepository"`
}

var roleUserService *RoleUserService = &RoleUserService{}

func init() {
	inject.InjectValue("roleUserService", roleUserService)
}

func (u *RoleUserService) ListMembers(roleId string, userName string, pageable query.Pageable) (int64, []domain.RoleMember) {
	return u.RoleUserRepository.ListMembers(roleId, userName, pageable)
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
