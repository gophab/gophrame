package auth

// 菜单分配文件相关的数据类型
type AuthNode struct {
	Id       int64      `json:"id" primaryKey:"yes"`
	Fid      int64      `json:"fid" fid:"Id"`
	Title    string     `json:"title,omitempty"`
	NodeType string     `json:"nodeType"`
	Expand   bool       `json:"expand"`
	Sort     int        `json:"sort"`
	Children []AuthNode `gorm:"-" json:"children,omitempty"`
}
