package auth

import "github.com/wjshen/gophrame/domain"

type RoleMenu struct {
	domain.Relation
	RoleId int64  `form:"role_id" gorm:"column:role_id" json:"roleId"`
	MenuId int64  `form:"menu_id" gorm:"column:menu_id" json:"menuId"`
	Status int    `form:"status" json:"status"`
	Remark string `form:"remark" json:"remark,omitempty"`
}

// 表名
func (a *RoleMenu) TableName() string {
	return "auth_role_menu"
}
