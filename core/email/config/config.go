package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
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
}

var Setting *EmailSetting = &EmailSetting{
	Enabled: false,
	Sender: struct {
	}{},
}

func init() {
	logger.Debug("Register Email Config")
	config.RegisterConfig("email", Setting, "Email Settings")
}
