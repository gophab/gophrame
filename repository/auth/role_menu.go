package auth

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/domain/auth"

	"gorm.io/gorm"
)

type RoleMenuRepository struct {
	*gorm.DB `inject:"database"`
}

var roleMenuRepository *RoleMenuRepository = &RoleMenuRepository{}

func init() {
	inject.InjectValue("roleMenuRepository", roleMenuRepository)
}

// 根据id获取菜单id
func (a *RoleMenuRepository) GetByRoleIds(ids []string) (result []auth.RoleMenu) {
	a.Model(&auth.RoleMenu{}).Where("role_id IN ?", ids).Select("distinct menu_id").Find(&result)
	return
}

// 根据postID获取菜单ID
func (a *RoleMenuRepository) GetByRoleId(id string) (result []auth.RoleMenu) {
	a.Model(&auth.RoleMenu{}).Where("role_id = ?", id).Select("distinct menu_id").Find(&result)
	return
}
