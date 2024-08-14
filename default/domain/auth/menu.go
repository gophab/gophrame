package auth

import (
	"github.com/gophab/gophrame/domain"
)

type MenuInfo struct {
	OperationInfo
	domain.ParametersEnabled
	Path      string `gorm:"column:path" json:"path,omitempty"`
	Component string `gorm:"column:component" json:"component,omitempty"`
	OutPage   bool   `gorm:"column:out_page" json:"outPage"`
	Hidden    bool   `gorm:"column:hidden;default:0" json:"hidden"`
}

type Menu struct {
	domain.AuditingEnabled
	MenuInfo
	Leaf     bool      `gorm:"->;default:true" json:"leaf"`
	Children []*Menu   `gorm:"-" json:"children,omitempty"`
	Buttons  []*Button `gorm:"foreignkey:Fid" json:"buttons,omitempty"`
}

// 表名
func (a *Menu) TableName() string {
	return "auth_menu"
}

// 系统菜单以及子表数据结构
type MenuCreate struct {
	MenuInfo
	ButtonArray  `json:"button_array"`
	ButtonDelete string `json:"button_delete"`
}

// 菜单主表以及子表修改的数据结构
type MenuEdit struct {
	MenuInfo
	ButtonArray  `json:"button_array"`
	ButtonDelete string `json:"button_delete"`
	Id           int64  `json:"id"`
}

// 数据类型被使用时，shouldbindjson 对于数字是可以接受  int  int64   float64, shouldbind 函数对于数字只能接受  float64
type ButtonArray []struct {
	MenuId           int64  `json:"menu_id"`
	ButtonId         int64  `json:"button_id"`
	RequestUrl       string `json:"request_url"`
	RequestMethod    string `json:"request_method"`
	Remark           string `json:"remark,omitempty"`
	Status           int64  `json:"status"`
	CreatedTime      string
	LastModifiedTime string
}
