package dto

import (
	"time"

	"github.com/gophab/gophrame/core/mapper"
	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/service/dto"
)

type User struct {
	dto.User
	CreatedTime      *time.Time `json:"createdTime"`
	LastModifiedTime *time.Time `json:"lastModifiedTime"`
	DelFlag          *bool      `json:"del_flag"`
	Admin            *bool      `json:"admin"`
	Status           *int       `json:"status"`
	Avatar           *string    `json:"avatar,omitempty"`
	Remark           *string    `json:"remark,omitempty"`
	LoginTimes       *int       `json:"loginTimes"`
	LastLoginTime    *time.Time `json:"lastLoginTime,omitempty"`
	LastLoginIp      *string    `json:"lastLoginIp,omitempty"`
	Roles            []Role     `json:"roles,omitempty"`
}

func (u *User) AsDomain() *domain.User {
	var result = domain.User{}
	if err := mapper.Map(u, &result); err != nil {
		return nil
	}
	return &result
}

func (a *User) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["del_flag"] = false
	return maps
}
