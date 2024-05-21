package domain

type Role struct {
	Entity
	Name string `json:"name"`
}

func (*Role) TableName() string {
	return "sys_role"
}
