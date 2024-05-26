package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type SocialSetting struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

var Setting *SocialSetting = &SocialSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Social Config")
	config.RegisterConfig("social", Setting, "Social Settings")
}
