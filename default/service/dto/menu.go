package dto

type Menu struct {
	Id     string `json:"id"`
	Name   string
	Path   string
	Type   string
	Method string

	CreatedBy  string
	ModifiedBy string
}

func (a *Menu) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["del_flag"] = false
	return maps
}
