package dto

import "time"

type SocialUser struct {
	Id            *string
	UserId        *string
	Name          *string
	Mobile        *string
	Email         *string
	Avatar        *string
	LoginTimes    int
	LastLoginTime *time.Time
	LastLoginIp   *string
}
