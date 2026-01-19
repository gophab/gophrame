package dto

type RoleCreate struct {
	Name   string `json:"name,omitempty"`
	Status int    `json:"status,omitempty"`
	Remark string `json:"remark,omitempty"`
}

type Role struct {
	RoleCreate
	Id string `json:"id,omitempty"`
	// CreatedBy      string `json:"createdBy,omitempty`
	// LastModifiedBy string `json:"lastModifiedBy,omitempty`
}

func (a *Role) GetMaps() map[string]any {
	maps := make(map[string]any)
	maps["del_flag"] = false
	return maps
}
