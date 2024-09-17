package domain

import "github.com/gophab/gophrame/domain"

type ModuleInfo struct {
	Name        string `gorm:"column:name" json:"name"`
	Title       string `gorm:"column:title" json:"title,omitempty" i18n:"yes"`
	Description string `gorm:"column:description" json:"description,omitempty" i18n:"yes"`
	Type        string `gorm:"column:type;default:MODULE" json:"type"`
	Oid         int64  `gorm:"column:oid;default:null" json:"oid,omitempty"`
	Fid         int64  `gorm:"column:fid;default:null" json:"fid,omitempty"`
	Fids        string `gorm:"column:fids" json:"fids"`
	Path        string `gorm:"column:path" json:"path"`
	Status      int    `gorm:"column:status;default:1" json:"status"`
}

type Module struct {
	domain.AuditingModel
	ModuleInfo

	Leaf     bool      `gorm:"-" json:"leaf"`
	Children []*Module `gorm:"-" json:"children,omitempty"`
}

func (*Module) TableName() string {
	return "auth_module"
}
