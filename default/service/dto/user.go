package dto

import (
	"time"

	"github.com/wjshen/gophrame/core/mapper"
	"github.com/wjshen/gophrame/default/domain"
)

type User struct {
	Id               *string    `json:"id"`
	CreatedTime      *time.Time `json:"createdTime"`
	LastModifiedTime *time.Time `json:"lastModifiedTime"`
	TenantId         *string    `json:"tenantId"`
	DelFlag          *bool      `json:"del_flag"`
	Login            *string    `json:"login,omitempty"`
	Mobile           *string    `json:"mobile,omitempty"`
	Email            *string    `json:"email,omitempty"`
	Name             *string    `json:"name,omitempty"`
	Status           *int       `json:"status"`
	Avatar           *string    `json:"avatar,omitempty"`
	Remark           *string    `json:"remark,omitempty"`
	LoginTimes       *int       `json:"loginTimes"`
	LastLoginTime    *time.Time `json:"lastLoginTime,omitempty"`
	LastLoginIp      *string    `json:"lastLoginIp,omitempty"`
	Password         *string    `json:"-"`
	Admin            *bool      `json:"admin"`
	InviterId        *string    `json:"inviterId,omitempty"`
	InviteCode       string     `json:"inviteCode,omitempty"`
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
