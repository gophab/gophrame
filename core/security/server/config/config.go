package config

import "time"

type OAuth2ServerSetting struct {
	Enabled                bool          `json:"enabled" yaml:"enabled"`
	AccessTokenExpireTime  time.Duration `json:"access_token_expire_time" yaml:"accessTokenExpireTime"`
	RefreshTokenExpireTime time.Duration `json:"refresh_token_expire_time" yaml:"refreshTokenExpireTime"`
}

var Setting *OAuth2ServerSetting = &OAuth2ServerSetting{
	Enabled:                false,
	AccessTokenExpireTime:  time.Hour * 8,
	RefreshTokenExpireTime: time.Hour * 24 * 100,
}
