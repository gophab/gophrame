package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type I18nSetting struct {
	// Common Settings
	Enabled bool `yaml:"enabled" json:"enabled"`
}

var Setting *I18nSetting = &I18nSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register I18n Config")
	config.RegisterConfig("i18n", Setting, "I18n Settings")
}
