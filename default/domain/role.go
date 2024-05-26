package domain

import "github.com/gophab/gophrame/domain"

type Role struct {
	domain.Entity
	Name string `json:"name"`
}

func (*Role) TableName() string {
	return "sys_role"
}
