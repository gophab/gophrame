package auth

import "github.com/wjshen/gophrame/domain"

// 用户在指定页面已分配的按钮列表
type Button struct {
	domain.Model
	CnName      string `json:"cnName,omitempty"`
	EnName      string `json:"enName,omitempty"`
	Color       string `json:"color,omitempty"`
	AllowMethod string `json:"allowMethod,omitempty"`
	Status      int    `json:"status"`
	Remark      string `json:"remark,omitempty"`
}

// 表名
func (b *Button) TableName() string {
	return "auth_button"
}
