package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/security/token/jwt"
)

type RedisTokenStoreSetting struct {
	Database  int    `json:"database" yaml:"database"`
	KeyPrefix string `json:"keyPrefix" yaml:"keyPrefix"`
}

type DatabaseTokeStoreSetting struct {
	Database    string `json:"database" yaml:"database"`
	TablePrefix string `json:"tablePrefix" yaml:"tablePrefix"`
}

type FileTokeStoreSetting struct {
	FileName string `json:"fileName" yaml:"fileName"`
}

type TokeStoreSetting struct {
	Mode     string                    `json:"mode" yaml:"mode"`
	Redis    *RedisTokenStoreSetting   `json:"redis" yaml:"redis"`
	Database *DatabaseTokeStoreSetting `json:"database" yaml:"database"`
	File     *FileTokeStoreSetting     `json:"file" yaml:"file"`
}

type TokenSetting struct {
	BindContextKey         string            `json:"bindContextKey"`
	HeaderTokenKey         string            `json:"headerTokenKey"`
	Cache                  bool              `json:"cache"`
	Store                  *TokeStoreSetting `json:"store" yaml:"store"`
	OnlineUsers            int               `json:"onlineUsers" yaml:"onlineUsers"`
	ReuseAccessToken       bool              `json:"reuseAccessToken" yaml:"reuseAccessToken"`
	ReuseRefreshToken      bool              `json:"reuseRefreshToken" yaml:"reuseRefreshToken"`
	AccessTokenExpireTime  time.Duration     `json:"accessTokenExpireTime" yaml:"accessTokenExpireTime"`
	RefreshTokenExpireTime time.Duration     `json:"refreshTokenExpireTime" yaml:"refreshTokenExpireTime"`
	UseJwtToken            bool              `json:"useJwtToken"`
	Jwt                    *jwt.JwtSetting   `json:"jwt" yaml:"jwt"`
}

var Setting *TokenSetting = &TokenSetting{
	OnlineUsers:            10,
	ReuseAccessToken:       true,
	ReuseRefreshToken:      true,
	AccessTokenExpireTime:  time.Hour * 8,
	RefreshTokenExpireTime: time.Hour * 24 * 100,

	// TokenStore配置
	Store: &TokeStoreSetting{
		Mode: "default",
	},
}

func init() {
	config.RegisterConfig("security.token", Setting, "Token Settings")
}
