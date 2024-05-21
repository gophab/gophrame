package dto

import (
	"github.com/wjshen/gophrame/domain"
)

type User struct {
	domain.User
}

func (a *User) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["del_flag"] = false
	return maps
}
