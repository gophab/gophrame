package auth

import (
	"github.com/gophab/gophrame/domain"
)

// 菜单分配文件相关的数据类型
// Type: [menu | button]
type OperationInfo struct {
	domain.PropertiesEnabled
	Id     int64  `gorm:"primaryKey" json:"id" primaryKey:"yes"`
	Fid    int64  `gorm:"column:fid;default:0" json:"fid" fid:"Id"`
	Name   string `gorm:"column:name" json:"name"`
	Icon   string `gorm:"column:icon" json:"icon,omitempty"`
	Title  string `gorm:"column:title" json:"title,omitempty" i18n:"yes"`
	Type   string `gorm:"column:type;->" json:"type"`
	Tags   string `gorm:"column:tags" json:"tags,omitempty"`
	Sort   int    `gorm:"column:sort;default:100" json:"sort"`
	Status int    `gorm:"column:status;default:1" json:"status"`
	Remark string `gorm:"column:remark" json:"remark,omitempty"`
	Level  int    `gorm:"column:level;default:0" json:"level"`
	Expand bool   `gorm:"column:expand;->" json:"expand"`
	Leaf   bool   `gorm:"column:leaf;->" json:"leaf"`
	Checks string `gorm:"column:cks;->" json:"checks"`
}

type Operation struct {
	OperationInfo
	Children []*Operation `gorm:"-" json:"children,omitempty"`
}
