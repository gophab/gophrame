package auth

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/service"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthRepository "github.com/gophab/gophrame/default/repository/auth"
	DefaultService "github.com/gophab/gophrame/default/service"
)

type AuthorityService struct {
	service.BaseService
	AuthorityRepository  *AuthRepository.AuthorityRepository  `inject:"authorityRepository"`
	MenuRepositoy        *AuthRepository.MenuRepository       `inject:"menuRepository"`
	MenuButtonRepository *AuthRepository.MenuButtonRepository `inject:"menuButtonRepository"`
	RoleUserService      *DefaultService.RoleUserService      `inject:"roleUserService"`
}

var authorityService *AuthorityService = &AuthorityService{}

func init() {
	inject.InjectValue("authorityService", authorityService)
}

func (u *AuthorityService) GetUserMenus(userId string) []AuthModel.Menu {
	roleIds := u.RoleUserService.GetUserRoleIds(userId)

	//根据岗位ID获取拥有的菜单ID,去重
	roleMenus := roleMenuService.GetByRoleIds(roleIds)
	menuIds := []int64{}
	for _, v := range roleMenus {
		menuIds = append(menuIds, v.MenuId)
	}

	//根据菜单 Ids数组 获取菜单信息
	return u.MenuRepositoy.GetByIds(menuIds)

}

func (u *AuthorityService) GetUserMenuTree(userId string) []AuthModel.Menu {
	menus := u.GetUserMenus(userId)
	if len(menus) > 1 {
		var dest = make([]AuthModel.Menu, 0)
		if err := util.CreateSqlResFormatFactory().ScanToTreeData(menus, &dest); err == nil {
			return dest
		} else {
			logger.Error("根据用户id查询权限范围内的菜单数据树形化出错", err.Error())
		}
	}
	return []AuthModel.Menu{}
}

// 查询用户打开指定的页面所拥有的按钮列表
func (u *AuthorityService) GetButtonListByMenuId(userId string, menuId int64) []AuthModel.Button {
	roleIds := u.RoleUserService.GetUserRoleIds(userId)
	if list := u.AuthorityRepository.GetButtonListByMenuId(roleIds, menuId); len(list) > 0 {
		return list
	}
	return []AuthModel.Button{}
}

func (u *AuthorityService) GetSystemAuthorities() (int64, []AuthModel.AuthNode) {
	return u.AuthorityRepository.GetSystemMenuButtonList()
}

func (u *AuthorityService) GetRoleAuthorities(roleId string) (int64, []AuthModel.AuthNode) {
	return u.AuthorityRepository.GetAssignedMenuButtonList(roleId)
}

// 根据用户ID获取所有权限的来源
func (a *AuthorityService) GetUserAuthorities(userId string) []AuthModel.AuthNode {
	//根据用户ID,查询隶属哪些组织机构
	result := a.AuthorityRepository.GetUserAuthorities(userId)
	if len(result) == 0 {
		return []AuthModel.AuthNode{}
	}
	return result
}
