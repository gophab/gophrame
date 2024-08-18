package auth

import (
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/service"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthRepository "github.com/gophab/gophrame/default/repository/auth"

	DefaultService "github.com/gophab/gophrame/default/service"
)

type AuthorityService struct {
	service.BaseService
	AuthorityRepository  *AuthRepository.AuthorityRepository  `inject:"authorityRepository"`
	MenuRepository       *AuthRepository.MenuRepository       `inject:"menuRepository"`
	MenuButtonRepository *AuthRepository.MenuButtonRepository `inject:"menuButtonRepository"`
	MenuService          *MenuService                         `inject:"menuService"`
	RoleUserService      *DefaultService.RoleUserService      `inject:"roleUserService"`
}

var authorityService *AuthorityService = &AuthorityService{}

func init() {
	inject.InjectValue("authorityService", authorityService)
}

func (u *AuthorityService) GetUserMenus(userId string) []*AuthModel.Menu {
	roleIds := u.RoleUserService.GetUserRoleIds(userId)

	//根据岗位ID获取拥有的菜单ID,去重
	menus, _ := u.AuthorityRepository.GetMenuByRoleIds(roleIds)

	//根据菜单 Ids数组 获取菜单信息
	return menus
}

func (u *AuthorityService) GetUserMenuTree(userId string) []*AuthModel.Menu {
	menus := u.GetUserMenus(userId)
	if len(menus) > 1 {
		return u.MenuService.MakeTree(menus)
	}
	return []*AuthModel.Menu{}
}

// 查询用户打开指定的页面所拥有的按钮列表
func (u *AuthorityService) GetButtonListByMenuId(userId string, menuId int64) []*AuthModel.Button {
	roleIds := u.RoleUserService.GetUserRoleIds(userId)
	if list, err := u.AuthorityRepository.GetButtonListByMenuId(roleIds, menuId); err == nil {
		return list
	}
	return []*AuthModel.Button{}
}

func (u *AuthorityService) GetSystemAuthorities() (int64, []*AuthModel.Operation) {
	return u.AuthorityRepository.GetSystemAuthorities()
}

func (u *AuthorityService) SetRoleAuthorities(roleId string, authorities []*AuthModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByRoleId(roleId, authorities)
}

func (u *AuthorityService) SetUserAuthorities(userId string, authorities []*AuthModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByUserId(userId, authorities)
}

func (u *AuthorityService) SetOrganizationAuthorities(organizationId string, authorities []*AuthModel.Operation) {
	u.AuthorityRepository.SetAuthoritiesByOrganizationId(organizationId, authorities)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAuthorities(userId string) []*AuthModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	_, result := a.AuthorityRepository.GetAuthoritiesByUserId(userId)
	if len(result) == 0 {
		return []*AuthModel.Operation{}
	}
	return result
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAvailableAuthorities(userId string) []*AuthModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetUserAuthorities(userId)
	if len(result) == 0 {
		return []*AuthModel.Operation{}
	}
	return result
}

func (u *AuthorityService) GetRoleAuthorities(roleId string) (int64, []*AuthModel.Operation) {
	return u.AuthorityRepository.GetAuthoritiesByRoleId(roleId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetRoleAvailableAuthorities(roleId string) []*AuthModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetRoleAuthorities(roleId)
	if len(result) == 0 {
		return []*AuthModel.Operation{}
	}
	return result
}

func (u *AuthorityService) GetOrganizationAuthorities(organizationId string) (int64, []*AuthModel.Operation) {
	return u.AuthorityRepository.GetAuthoritiesByOrganizationId(organizationId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetOrganizationAvailableAuthorities(organizationId string) []*AuthModel.Operation {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetOrganizationAuthorities(organizationId)
	if len(result) == 0 {
		return []*AuthModel.Operation{}
	}
	return result
}

func (a *AuthorityService) makeTree(src []*AuthModel.Operation, dest *[]*AuthModel.Operation) error {
	var result = *dest
	var srcMap = make(map[int64]*AuthModel.Operation)
	for _, item := range src {
		srcMap[item.Id] = item
	}
	for _, item := range src {
		if item.Fid != 0 {
			var parent = srcMap[item.Fid]
			if parent != nil {
				if parent.Children == nil {
					parent.Children = make([]*AuthModel.Operation, 0)
				}
				parent.Children = append(parent.Children, item)
			}
		} else {
			result = append(result, item)
		}
	}
	*dest = result
	return nil
}
