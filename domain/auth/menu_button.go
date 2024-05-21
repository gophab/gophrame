package auth

import "github.com/wjshen/gophrame/domain"

type MenuButtonRelation struct {
	domain.Relation
	MenuId        int64  `json:"menuId"`
	ButtonId      int64  `json:"buttonId"`
	RequestMethod string `json:"requestMethod,omitempty"`
	RequestUrl    string `json:"requestUrl,omitempty"`
	Status        int    `json:"status"`
	Remark        string `json:"remark,omitempty"`
}

type MenuButton struct {
	MenuButtonRelation
	Id          int64        `gorm:"->" json:"id" primaryKey:"yes"`
	Fid         int64        `gorm:"->" json:"fid" fid:"Id"`
	Icon        string       `gorm:"->" json:"icon,omitempty"`
	Name        string       `gorm:"->" json:"name,omitempty"`
	Loading     bool         `gorm:"->" json:"loading"`
	Path        string       `gorm:"->" json:"path,omitempty"`
	Component   string       `gorm:"->" json:"component,omitempty"`
	NodeLevel   int          `gorm:"->" json:"nodeLevel"`
	OutPage     bool         `gorm:"->" json:"outPage"`
	Sort        int          `gorm:"->" json:"sort"`
	Title       string       `gorm:"->" json:"title,omitempty"`
	MenuFid     int64        `gorm:"->" json:"menuFid,omitempty"`
	ButtonName  string       `gorm:"->" json:"buttonName,omitempty"`
	ButtonColor string       `gorm:"->" json:"buttonColor,omitempty"`
	NodeType    string       `gorm:"->" json:"nodeType,omitempty"`
	Expand      int8         `gorm:"->" json:"expand"`
	Children    []MenuButton `gorm:"-" json:"children,omitempty"`
}

// 表名
func (a *MenuButton) TableName() string {
	return "auth_menu_button"
}
