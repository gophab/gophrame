package domain

import (
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/domain"

	"gorm.io/gorm"
)

type RoleUser struct {
	domain.Relation
	domain.TenantEnabled
	RoleId string `gorm:"column:role_id;primaryKey" json:"roleId"`
	UserId string `gorm:"column:user_id;primaryKey" json:"userId"`
	Status int    `gorm:"column:status;default:1" json:"status"`
	Remark string `gorm:"column:remark;default:null" json:"remark,omitempty"`
	Role   *Role  `gorm:"foreignKey:role_id" json:"role,omitempty"`
	User   *User  `gorm:"foreignKey:user_id" json:"user,omitempty"`
}

func (m *RoleUser) BeforeCreate(tx *gorm.DB) (err error) {
	if m.TenantId == "" {
		m.TenantId = SecurityUtil.GetCurrentTenantId(nil)
	}

	m.TenantEnabled.BeforeCreate(tx)
	return m.Relation.BeforeCreate(tx)
}

func (m *RoleUser) BeforeSave(tx *gorm.DB) (err error) {
	return m.Relation.BeforeSave(tx)
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
