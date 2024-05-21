package domain

import "time"

type SocialUser struct {
	Entity
	DelFlag       bool       `gorm:"column:del_flag" json:"delFlag"`
	Type          string     `gorm:"column:type" json:"type"`
	OpenId        *string    `gorm:"column:open_id" json:"openId,omitempty"`
	SocialId      *string    `gorm:"column:social_id" json:"socialId,omitempty"`
	Mobile        *string    `gorm:"column:mobile" json:"mobile,omitempty"`
	Email         *string    `gorm:"column:email" json:"email,omitempty"`
	Name          *string    `gorm:"column:name" json:"name,omitempty"`
	Status        *int       `gorm:"column:status;default:1" json:"status"`
	Avatar        *string    `gorm:"column:avatar" json:"avatar,omitempty"`
	LoginTimes    int        `gorm:"column:login_times" json:"loginTimes"`
	Remark        *string    `gorm:"column:remark" json:"remark,omitempty"`
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"lastLoginTime,omitempty"`
	LastLoginIp   *string    `gorm:"column:last_login_ip" json:"lastLoginIp,omitempty"`
	NickName      *string    `gorm:"column:nick_name" json:"nickName,omitempty"`
	Title         *string    `gorm:"column:title;" json:"title,omitempty"`
	UserId        *string    `gorm:"column:user_id" json:"userId,omitempty"`
	Roles         []Role     `gorm:"-" json:"roles,omitempty"`
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
