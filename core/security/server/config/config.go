package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
)

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

func init() {
	config.RegisterConfig("security.server", Setting, "OAuth2 Server Settings")
}
