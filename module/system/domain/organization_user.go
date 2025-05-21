package domain

import "github.com/gophab/gophrame/domain"

type OrganizationUser struct {
	domain.Relation
	OrganizationId string `gorm:"column:organization_id" json:"organizationId"`
	UserId         string `gorm:"column:user_id" json:"userId"`
	PositionId     *int64 `gorm:"column:position_id" json:"positionId"`
	Status         int    `gorm:"column:status;default:1" json:"status"`
	Remark         string `json:"remark,omitempty"`
}

// 表名
func (a *OrganizationUser) TableName() string {
	return "sys_organization_user"
}

// 定义不同的查询结果返回的数据结构体
type OrganizationMember struct {
	OrganizationUser
	Name             string `gorm:"->" json:"name,omitempty"`
	Login            string `gorm:"->" json:"login,omitempty"`
	Mobile           string `gorm:"->" json:"mobile,omitempty"`
	Email            string `gorm:"->" json:"email,omitempty"`
	Title            string `gorm:"->" json:"title,omitempty"`
	OrganizationName string `gorm:"->" json:"organizationName,omitempty"`
	PositionName     string `gorm:"->" json:"positionName,omitempty"`
}
