package domain

import (
	"time"

	"github.com/gophab/gophrame/domain"
	"gorm.io/gorm"

	"github.com/gophab/gophrame/core/util"
)

type UserInfo struct {
	Login         *string    `gorm:"column:login" json:"login,omitempty"`
	Mobile        *string    `gorm:"column:mobile" json:"mobile,omitempty"`
	Email         *string    `gorm:"column:email" json:"email,omitempty"`
	Name          *string    `gorm:"column:name" json:"name,omitempty" i18n:"yes"`
	Status        *int       `gorm:"column:status" json:"status"`
	Avatar        *string    `gorm:"column:avatar" json:"avatar,omitempty"`
	Remark        *string    `gorm:"column:remark" json:"remark,omitempty"`
	InviterId     *string    `gorm:"column:inviter_id" json:"inviterId,omitempty"`
	LoginTimes    int        `gorm:"column:login_times;default:0" json:"loginTimes"`
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"lastLoginTime,omitempty"`
	LastLoginIp   *string    `gorm:"column:last_login_ip" json:"lastLoginIp,omitempty"`
}

type User struct {
	domain.DeletableEntity
	UserInfo
	Password         string  `gorm:"column:password" json:"-"`
	Admin            bool    `gorm:"column:admin;default:false" json:"admin"`
	Roles            []*Role `gorm:"many2many:sys_role_user;" json:"roles,omitempty"`        /* 角色 */
	InviteCode       string  `gorm:"-" json:"inviteCode,omitempty"`                          /* 邀请码 */
	OrganizationId   *int64  `gorm:"column:organization_id" json:"organizationId,omitempty"` /* 所在组织ID */
	OrganizationName string  `gorm:"-" json:"organizationName,omitempty"`                    /* 所在组织 */
	Tenant           *Tenant `gorm:"-" json:"tenant,omitempty"`                              /* 所在企业 */
}

func (e *User) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Mobile != nil {
		e.Mobile = util.StringAddr(util.FullPhoneNumber(util.NotNullString(e.Mobile)))
	}
	// e.DeleteEnabled.BeforeCreate(tx)
	return e.DeletableEntity.BeforeCreate(tx)
}

func (e *User) BeforeSave(tx *gorm.DB) (err error) {
	if e.Mobile != nil {
		e.Mobile = util.StringAddr(util.FullPhoneNumber(util.NotNullString(e.Mobile)))
	}
	// e.Entity.BeforeSave(tx)
	return e.DeletableEntity.BeforeSave(tx)
}

func (u *User) TableName() string {
	return "sys_user"
}

func (u *User) SetPassword(value string) *User {
	u.Password = util.SHA1(value)
	return u
}

func (u *User) HasRole(role string) bool {
	if len(u.Roles) > 0 {
		for _, r := range u.Roles {
			if r.Id == role || r.Name == role {
				return true
			}
		}
	}

	return false
}
