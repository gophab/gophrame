package social

import (
	"context"
	"time"
)

type SocialUser struct {
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
}

func (u *SocialUser) SetSocialId(socialType string, socialId string) *SocialUser {
	u.Type = socialType
	u.SocialId = &socialId
	return u
}

func (u *SocialUser) GetId() string {
	if u.SocialId != nil {
		return u.Type + "_" + *u.SocialId
	}
	return ""
}

type SocialService interface {
	GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *SocialUser
}

type SocialManager struct {
	Services map[string]SocialService
}

func GetSocialService(social string) SocialService {
	return Manager.GetSocialService(social)
}

func RegisterSocialService(social string, service SocialService) error {
	return Manager.RegisterSocialService(social, service)
}

var Manager = &SocialManager{
	Services: make(map[string]SocialService, 10),
}

func (m *SocialManager) GetSocialService(social string) SocialService {
	if service, b := m.Services[social]; b {
		return service
	}
	return nil
}

func (m *SocialManager) RegisterSocialService(social string, service SocialService) error {
	m.Services[social] = service
	return nil
}
