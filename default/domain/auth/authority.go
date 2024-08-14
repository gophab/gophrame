package auth

import "github.com/gophab/gophrame/domain"

type Authority struct {
	AuthType string `gorm:"column:auth_type;primaryKey" json:"authType"`
	AuthId   int64  `gorm:"column:auth_id;primaryKey" json:"authId"`
	Sort     int    `gorm:"column:sort;default:100" json:"sort"`
	Status   int    `gorm:"column:status;default:1" json:"status"`
}

type RoleAuthority struct {
	domain.AuditingEnabled
	Authority
	RoleId string `gorm:"column:role_id;primaryKey" json:"roleId"`
}

func (*RoleAuthority) TableName() string {
	return "auth_role_authority"
}

type UserAuthority struct {
	domain.AuditingEnabled
	Authority
	UserId string `gorm:"column:user_id;primaryKey" json:"userId"`
}

func (*UserAuthority) TableName() string {
	return "auth_user_authority"
}

type OrganizationAuthority struct {
	domain.AuditingEnabled
	Authority
	OrganizationId string `gorm:"column:organization_id;primaryKey" json:"organizationId"`
}

func (*OrganizationAuthority) TableName() string {
	return "auth_organization_authority"
}
