package domain

import "github.com/gophab/gophrame/domain"

type Organization struct {
	domain.Model
	Fid      int64          `json:"fid" fid:"Id"`
	Name     string         `json:"name"`
	Status   string         `json:"status"`
	PathInfo string         `json:"pathInfo"`
	Remark   string         `json:"remark,omitempty"`
	Leaf     bool           `gorm:"->" json:"leaf"` // 是否为叶子节点
	Expand   bool           `gorm:"->" json:"expand"`
	NodeType string         `gorm:"->" json:"nodeType"`
	Children []Organization `gorm:"-" json:"children,omitempty"`
}

func (*Organization) TableName() string {
	return "sys_organization"
}
