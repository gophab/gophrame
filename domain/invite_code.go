package domain

import "time"

type InviteCode struct {
	Entity
	InviteCode   string     `gorm:"column:invite_code" json:"inviteCode"`
	UserId       string     `gorm:"column:user_id" json:"userId"`
	Channel      string     `gorm:"column:channel" json:"channel"`
	ExpireTime   *time.Time `gorm:"column:expire_time" json:"expireTime,omitempty"`
	InviteLimit  int64      `gorm:"column:invite_limit" json:"inviteLimit"`
	InvitedLimit int64      `gorm:"column:invited_limit" json:"invitedLimit"`
}

func (*InviteCode) TableName() string {
	return "sys_invite_code"
}

func (c *InviteCode) IsExpired() bool {
	if c.ExpireTime == nil {
		return false
	}

	if c.ExpireTime.After(time.Now()) {
		return false
	}

	return true
}
