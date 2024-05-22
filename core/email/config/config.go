package config

import (
	CodeConfig "github.com/wjshen/gophrame/core/code/config"
	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"

	EmailCodeConfig "github.com/wjshen/gophrame/core/email/code/config"
)

type RedisCodeStoreSetting struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Database  int    `json:"database" yaml:"database"`
	KeyPrefix string `json:"keyPrefix" yaml:"keyPrefix"`
}

type CacheCodeStoreSetting struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

type EmailSetting struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
	Sender  struct {
	}
	Store *CodeConfig.CodeStoreSetting `json:"store" yaml:"store"`
}

var Setting *EmailSetting = &EmailSetting{
	Enabled: false,
	Sender: struct {
	}{},
	Store: EmailCodeConfig.Setting,
}

func init() {
	logger.Debug("Register Email Config")
	config.RegisterConfig("email", Setting, "Email Settings")
}
