package domain

import (
	"github.com/gophab/gophrame/domain"
)

type ContentTemplateInfo struct {
	domain.ParametersEnabled
	domain.PropertiesEnabled
	Name    string `gorm:"column:name" json:"name,omitempty" i18n:"true"`
	Title   string `gorm:"column:title" json:"title,omitempty" i18n:"true"`
	Type    string `gorm:"column:type" json:"type,omitempty"`
	Scene   string `gorm:"column:scene" json:"scene,omitempty"`
	Content string `gorm:"column:content" json:"content,omitempty" i18n:"true"`
	Status  int    `gorm:"column:status" json:"status"`
}

type ContentTemplate struct {
	domain.AuditingEntity
	ContentTemplateInfo
}

func (*ContentTemplate) TableName() string {
	return "sys_content_template"
}
