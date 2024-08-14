package domain

import "github.com/gophab/gophrame/domain"

type RoleInfo struct {
	Name        string  `gorm:"column:name" json:"name"`
	Title       string  `gorm:"column:title;default:null" json:"title,omitempty"`
	Description string  `gorm:"column:description;default:null" json:"description,omitempty"`
	Scope       string  `gorm:"column:scope;default:TENANT" json:"scope,omitempty"`
	Includes    string  `gorm:"column:includes" json:"includes,omitempty"`
	Children    []*Role `gorm:"-" json:"children,omitempty"`
}

type Role struct {
	domain.DeletableEntity
	RoleInfo
}

func (*Role) TableName() string {
	return "sys_role"
}
