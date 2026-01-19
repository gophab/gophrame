package domain

import (
	"github.com/gophab/gophrame/core/social"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/domain"
	"gorm.io/gorm"
)

type SocialUser struct {
	domain.DeletableEntity
	social.SocialUser
	Roles []Role `gorm:"-" json:"roles,omitempty"`
}

func (u *SocialUser) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = u.GetId()
	}

	if u.Mobile != nil {
		u.Mobile = util.StringAddr(util.FullPhoneNumber(util.NotNullString(u.Mobile)))
	}

	return u.DeletableEntity.BeforeCreate(tx)
}

func (u *SocialUser) BeforeSave(tx *gorm.DB) (err error) {
	if u.Mobile != nil {
		u.Mobile = util.StringAddr(util.FullPhoneNumber(util.NotNullString(u.Mobile)))
	}
	return u.DeletableEntity.BeforeSave(tx)
}

func (u *SocialUser) TableName() string {
	return "sys_social_user"
}

func (u *SocialUser) SetSocialId(socialType string, socialId string) *SocialUser {
	u.Type = socialType
	u.SocialId = &socialId
	u.Id = u.GetId()

	return u
}
