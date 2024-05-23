package domain

import (
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/domain"
)

type SocialUser struct {
	domain.Entity
	social.SocialUser
	Roles []Role `gorm:"-" json:"roles,omitempty"`
}

func (u *SocialUser) TableName() string {
	return "sys_social_user"
}

func (u *SocialUser) SetSocialId(socialType string, socialId string) *SocialUser {
	u.Type = socialType
	u.SocialId = &socialId
	u.Id = socialType + "_" + socialId

	return u
}
