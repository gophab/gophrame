package domain

import (
	"time"

	"github.com/wjshen/gophrame/core/util"
)

type UserBase struct {
	Entity
	DelFlag       bool       `gorm:"column:del_flag" json:"del_flag"`
	Login         *string    `gorm:"column:login" json:"login,omitempty"`
	Mobile        *string    `gorm:"column:mobile" json:"mobile,omitempty"`
	Email         *string    `gorm:"column:email" json:"email,omitempty"`
	Name          *string    `gorm:"column:name" json:"name,omitempty"`
	Status        *int       `gorm:"column:status" json:"status"`
	Avatar        *string    `gorm:"column:avatar" json:"avatar,omitempty"`
	Remark        *string    `gorm:"column:remark" json:"remark,omitempty"`
	InviterId     *string    `gorm:"column:inviter_id" json:"inviterId,omitempty"`
	LoginTimes    int        `gorm:"column:login_times" json:"loginTimes"`
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"lastLoginTime,omitempty"`
	LastLoginIp   *string    `gorm:"column:last_login_ip" json:"lastLoginIp,omitempty"`
}

type User struct {
	UserBase
	Password   string `gorm:"column:password" json:"-"`
	Admin      bool   `gorm:"column:admin" json:"admin"`
	InviteCode string `gorm:"-" json:"inviteCode,omitempty"`
	Roles      []Role `gorm:"many2many:sys_role_user;" json:"roles,omitempty"`
}

type UserWithOrganization struct {
	User
	OrganizationId   int64  `gorm:"->" json:"organizationId"`
	OrganizationName string `gorm:"->" json:"organizationName"`
}

func (u *User) TableName() string {
	return "sys_user"
}

func (u *User) SetPassword(value string) *User {
	u.Password = util.SHA1(value)
	return u
}
