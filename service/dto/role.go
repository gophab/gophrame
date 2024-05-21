package dto

type RoleCreate struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

type Role struct {
	RoleCreate
	Id         string `json:"id"`
	CreatedBy  string
	ModifiedBy string
}

func (a *Role) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["del_flag"] = false
	return maps
}
