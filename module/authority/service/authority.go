package service

import (
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/authority/repository"

	OperationModel "github.com/gophab/gophrame/module/operation/domain"
	Operation "github.com/gophab/gophrame/module/operation/service"
	System "github.com/gophab/gophrame/module/system/service"

	"github.com/gophab/gophrame/service"
)

type AuthorityService struct {
	service.BaseService
	AuthorityRepository *repository.AuthorityRepository `inject:"authorityRepository"`
	MenuService         *Operation.MenuService          `inject:"menuService"`
	RoleUserService     *System.RoleUserService         `inject:"roleUserService"`
}

var authorityService *AuthorityService = &AuthorityService{}

func init() {
	inject.InjectValue("authorityService", authorityService)
}

func (u *AuthorityService) GetUserMenus(userId string) []*OperationModel.Menu {
	roleIds := u.RoleUserService.GetUserRoleIds(userId)

	//根据岗位ID获取拥有的菜单ID,去重
	menus, _ := u.AuthorityRepository.GetMenuByRoleIds(roleIds)

	//根据菜单 Ids数组 获取菜单信息
	return menus
}

func (u *AuthorityService) GetUserMenuTree(userId string) []*OperationModel.Menu {
	menus := u.GetUserMenus(userId)
	if len(menus) > 1 {
		return u.MenuService.MakeTree(menus)
	}
	return []*OperationModel.Menu{}
}

func (u *AuthorityService) GetSystemAuthorities() (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetSystemAuthorities()
}

func (u *AuthorityService) SetRoleAuthorities(roleId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByRoleId(roleId, authorities)
}

func (u *AuthorityService) SetUserAuthorities(userId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByUserId(userId, authorities)
}

func (u *AuthorityService) SetOrganizationAuthorities(organizationId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByOrganizationId(organizationId, authorities)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAuthorities(userId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	_, result := a.AuthorityRepository.GetAuthoritiesByUserId(userId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAvailableAuthorities(userId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetUserAuthorities(userId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}

func (u *AuthorityService) GetRoleAuthorities(roleId string) (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetAuthoritiesByRoleId(roleId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetRoleAvailableAuthorities(roleId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetRoleAuthorities(roleId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}

func (u *AuthorityService) GetOrganizationAuthorities(organizationId string) (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetAuthoritiesByOrganizationId(organizationId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetOrganizationAvailableAuthorities(organizationId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetOrganizationAuthorities(organizationId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}
