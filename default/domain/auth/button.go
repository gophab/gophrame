package auth

import "github.com/gophab/gophrame/domain"

type ButtonInfo struct {
	OperationInfo
	Color       string `gorm:"column:color" json:"color,omitempty"`
	AllowMethod string `gorm:"column:allow_method" json:"allowMethod,omitempty"`
}

type Button struct {
	domain.AuditingEnabled
	ButtonInfo
}

// 表名
func (b *Button) TableName() string {
	return "auth_button"
}
