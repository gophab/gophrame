package domain

import (
	"github.com/gophab/gophrame/domain"
	"gorm.io/gorm"
)

type ButtonInfo struct {
	OperationInfo
	Color       string `gorm:"column:color" json:"color,omitempty"`
	AllowMethod string `gorm:"column:allow_method" json:"allowMethod,omitempty"`
}

type Button struct {
	domain.AuditingEnabled
	ButtonInfo
}

func (m *Button) BeforeCreate(tx *gorm.DB) (err error) {
	m.Entity.BeforeCreate(tx)
	return m.AuditingEnabled.BeforeCreate(tx)
}

// 表名
func (b *Button) TableName() string {
	return "auth_button"
}
