package domain

import "github.com/wjshen/gophrame/domain"

type RoleUser struct {
	domain.Relation
	RoleId string `gorm:"column:role_id" json:"roleId"`
	UserId string `gorm:"column:user_id" json:"userId"`
	Status int    `gorm:"column:status;default:1" json:"status"`
	Remark string `gorm:"column:remark" json:"remark,omitempty"`
}

// 表名
func (a *RoleUser) TableName() string {
	return "sys_role_user"
}

// 定义不同的查询结果返回的数据结构体
type RoleMember struct {
	RoleUser
	Name     string `json:"name,omitempty"`
	Login    string `json:"login,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Email    string `json:"email,omitempty"`
	RoleName string `json:"roleName,omitempty"`
}
