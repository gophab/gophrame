package service

import (
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/authority/domain"
	"github.com/gophab/gophrame/module/authority/repository"

	OperationModel "github.com/gophab/gophrame/module/operation/domain"
	Operation "github.com/gophab/gophrame/module/operation/service"
	System "github.com/gophab/gophrame/module/system/service"

	"github.com/gophab/gophrame/service"
)

type AuthorityService struct {
	service.BaseService
	AuthorityRepository     *repository.AuthorityRepository `inject:"authorityRepository"`
	MenuService             *Operation.MenuService          `inject:"menuService"`
	RoleUserService         *System.RoleUserService         `inject:"roleUserService"`
	OrganizationUserService *System.OrganizationUserService `inject:"organizationUserService"`
	OrganizationService     *System.OrganizationService     `inject:"organizationService"`
}

var authorityService *AuthorityService = &AuthorityService{}

func init() {
	inject.InjectValue("authorityService", authorityService)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAvailableAuthorities(userId string, authType string) []*domain.UserAuthority {
	var results = make([]*domain.UserAuthority, 0)
	// 1. 用户的组织授权
	var organizationIds = a.OrganizationUserService.GetUserOrganizationIds(userId)
	if len(organizationIds) > 0 {
		for _, organizationId := range organizationIds {
			if authorities, err := a.GetOrganizationAvailableAuthorities(organizationId, authType); err == nil && authorities != nil {
				for _, authority := range authorities {
					var userAuthority = &domain.UserAuthority{
						UserId: userId,
						Authority: domain.Authority{
							AuthType: authType,
							AuthId:   authority.AuthId,
							Status:   1,
						},
					}
					results = append(results, userAuthority)
				}
			}
		}
	}
	// 2. 用户的角色授权
	var roleIds = a.RoleUserService.GetUserRoleIds(userId)
	if len(roleIds) > 0 {
		for _, roleId := range roleIds {
			if authorities, err := a.GetRoleAuthorities(roleId, authType); err == nil && authorities != nil {
				for _, authority := range authorities {
					var userAuthority = &domain.UserAuthority{
						UserId: userId,
						Authority: domain.Authority{
							AuthType: authType,
							AuthId:   authority.AuthId,
							Status:   1,
						},
					}
					results = append(results, userAuthority)
				}
			}
		}
	}

	// 3. 用户授权
	userAuthorities, _ := a.AuthorityRepository.GetUserAuthorities(userId, authType)
	if len(userAuthorities) > 0 {
		results = append(results, userAuthorities...)
	}
	return results
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAuthorities(userId string, authType string) []*domain.UserAuthority {
	//根据用户ID,查询隶属哪些组织机构
	result, _ := a.AuthorityRepository.GetUserAuthorities(userId, authType)
	if len(result) == 0 {
		return []*domain.UserAuthority{}
	}
	return result
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAuthority(userId string, authType string, authId string) *domain.UserAuthority {
	//根据用户ID,查询隶属哪些组织机构
	if result, err := a.AuthorityRepository.GetUserAuthority(userId, authType, authId); err == nil {
		return result
	} else {
		return nil
	}
}

func (u *AuthorityService) GetRoleAuthorities(roleId string, authType string) ([]*domain.RoleAuthority, error) {
	return u.AuthorityRepository.GetRoleAuthorities(roleId, authType)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetRoleAuthority(roleId string, authType string, authId string) *domain.RoleAuthority {
	//根据用户ID,查询隶属哪些组织机构
	if result, err := a.AuthorityRepository.GetRoleAuthority(roleId, authType, authId); err == nil {
		return result
	} else {
		return nil
	}
}

func (u *AuthorityService) GetOrganizationAvailableAuthorities(organizationId string, authType string) ([]*domain.OrganizationAuthority, error) {
	organizations := u.OrganizationService.GetParentList(organizationId)
	if len(organizations) > 0 {
		var organizationIds = []string{organizationId}
		for _, organization := range organizations {
			organizationIds = append(organizationIds, organization.Id)
		}
		return u.AuthorityRepository.GetOrganizationsAuthorities(organizationIds, authType)
	} else {
		return u.AuthorityRepository.GetOrganizationAuthorities(organizationId, authType)
	}
}

func (u *AuthorityService) GetOrganizationAuthorities(organizationId string, authType string) ([]*domain.OrganizationAuthority, error) {
	return u.AuthorityRepository.GetOrganizationAuthorities(organizationId, authType)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetOrganizationAuthority(organizationId string, authType string, authId string) *domain.OrganizationAuthority {
	//根据用户ID,查询隶属哪些组织机构
	if result, err := a.AuthorityRepository.GetOrganizationAuthority(organizationId, authType, authId); err == nil {
		return result
	} else {
		return nil
	}
}

func (u *AuthorityService) SetRoleAuthorities(roleId string, authType string, authIds []string) {
	u.AuthorityRepository.SetAuthoritiesByRoleId(roleId, authType, authIds)
}

func (u *AuthorityService) SetUserAuthorities(userId string, authType string, authIds []string) {
	u.AuthorityRepository.SetAuthoritiesByUserId(userId, authType, authIds)
}

func (u *AuthorityService) SetOrganizationAuthorities(organizationId string, authType string, authIds []string) {
	u.AuthorityRepository.SetAuthoritiesByOrganizationId(organizationId, authType, authIds)
}

// Operations
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

func (u *AuthorityService) GetSystemOperations() (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetSystemOperations()
}

func (u *AuthorityService) SetRoleOperations(roleId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetOperationsByRoleId(roleId, authorities)
}

func (u *AuthorityService) SetUserOperations(userId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetOpertionsByUserId(userId, authorities)
}

func (u *AuthorityService) SetOrganizationOperations(organizationId string, authorities []*OperationModel.Operation) {
	u.AuthorityRepository.SetOperationsByOrganizationId(organizationId, authorities)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserOperations(userId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	_, result := a.AuthorityRepository.GetOpertionsByUserId(userId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAvailableOperations(userId string) []*OperationModel.Operation {
	var results = make([]*OperationModel.Operation, 0)
	// // 1. 用户的组织授权
	// var organizationIds = a.OrganizationUserService.GetUserOrganizationIds(userId)
	// if len(organizationIds) > 0 {
	// 	for _, organizationId := range organizationIds {
	// 		if operations := a.GetOrganizationAvailableOperations(organizationId); len(operations) > 0 {
	// 			results = append(results, operations...)
	// 		}
	// 	}
	// }
	// // 2. 用户的角色授权
	// var roleIds = a.RoleUserService.GetUserRoleIds(userId)
	// if len(roleIds) > 0 {
	// 	for _, roleId := range roleIds {
	// 		if _, operations := a.GetRoleOperations(roleId); len(operations) > 0 {
	// 			results = append(results, operations...)
	// 		}
	// 	}
	// }
	// // 3. 用户的授权
	// if operations := a.GetUserOperations(userId); len(operations) > 0 {
	// 	results = append(results, operations...)
	// }

	// 直接
	if operations := a.AuthorityRepository.GetUserAvailableOperations(userId); len(operations) > 0 {
		results = append(results, operations...)
	}

	return results
}

func (u *AuthorityService) GetRoleOperations(roleId string) (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetOperationsByRoleId(roleId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetRoleAvailableOperations(roleId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetRoleOperations(roleId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}

func (u *AuthorityService) GetOrganizationOperations(organizationId string) (int64, []*OperationModel.Operation) {
	return u.AuthorityRepository.GetOperationsByOrganizationId(organizationId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetOrganizationAvailableOperations(organizationId string) []*OperationModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetOrganizationOperations(organizationId)
	if len(result) == 0 {
		return []*OperationModel.Operation{}
	}
	return result
}
