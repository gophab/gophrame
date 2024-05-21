package auth

import (
	"time"

	"github.com/wjshen/gophrame/domain"
)

type BaseMenu struct {
	domain.Model
	Fid       int64  `json:"fid" fid:"Id"`
	Icon      string `json:"icon,omitempty"`
	Title     string `json:"title,omitempty"`
	Name      string `json:"name,omitempty"`
	Path      string `json:"path,omitempty"`
	Component string `json:"component,omitempty"`
	Status    int    `gorm:"column:status;default:1" json:"status"`
	OutPage   bool   `json:"outPage"`
	Sort      int    `json:"sort"`
	Remark    string `json:"remark,omitempty"`
	NodeLevel int    `json:"nodeLevel"`
}

type Menu struct {
	Id               int64     `gorm:"primaryKey" json:"id" primaryKey:"yes"`
	Fid              int64     `json:"fid" fid:"Id"`
	Icon             string    `json:"icon,omitempty"`
	Title            string    `json:"title,omitempty"`
	Name             string    `json:"name,omitempty"`
	Path             string    `json:"path,omitempty"`
	Component        string    `json:"component,omitempty"`
	Status           int       `json:"status"`
	OutPage          bool      `json:"outPage"`
	Sort             int       `json:"sort"`
	Remark           string    `json:"remark,omitempty"`
	NodeLevel        int       `json:"nodeLevel"`
	CreatedTime      time.Time `gorm:"autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"autoUpdateTime" json:"lastModifiedTime"`
	HasSubNode       bool      `gorm:"->" json:"hasSubNode"`
	Leaf             bool      `gorm:"->" json:"leaf"`
	Children         []Menu    `gorm:"-" json:"children,omitempty"`
	Loading          bool      `gorm:"-" json:"loading"`
}

// 表名
func (a *Menu) TableName() string {
	return "auth_menu"
}

type MenuWithButton struct {
	Id               int64     `gorm:"primaryKey" json:"id" primaryKey:"yes"`
	Fid              int64     `json:"fid"`
	Icon             string    `json:"icon,omitempty"`
	Title            string    `json:"title,omitempty"`
	Name             string    `json:"name,omitempty"`
	Path             string    `json:"path,omitempty"`
	Component        string    `json:"component,omitempty"`
	Status           int       `json:"status"`
	OutPage          bool      `json:"outPage"`
	Sort             int       `json:"sort"`
	Remark           string    `json:"remark,omitempty"`
	NodeLevel        int       `json:"nodeLevel"`
	CreatedTime      time.Time `gorm:"autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"autoUpdateTime" json:"lastModifiedTime"`
	Children         []struct {
		MenuId      int64  `fid:"Id" json:"menuId"`
		ButtonId    int64  `json:"buttonId" primaryKey:"yes"`
		ButtonName  string `json:"buttonName"`
		ButtonColor string `json:"buttonColor"`
	} `json:"buttons,omitempty" gorm:"-"`
}

// 系统菜单以及子表数据结构
type MenuCreate struct {
	Title        string `form:"title" json:"title" binding:"required,min=1"`
	Icon         string `form:"icon" json:"icon"`
	Fid          *int64 `form:"fid" json:"fid" binding:"required,min=0"`
	Status       *int   `form:"status" json:"status" binding:"required,min=0"`
	OutPage      *int   `form:"out_page" json:"out_page" binding:"required,min=0,max=1"`
	Sort         *int   `form:"sort" json:"sort" binding:"required,min=0"`
	Name         string `form:"name" json:"name" binding:"required,min=1"`
	Path         string `form:"path" json:"path" `
	Component    string `form:"component" json:"component"`
	Remark       string `form:"remark" json:"remark"`
	ButtonDelete string `json:"button_delete"`
	ButtonArray  `json:"button_array"`
}

// 菜单主表以及子表修改的数据结构
type MenuEdit struct {
	Id int64 `json:"id"`
	MenuCreate
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
