package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
)

type ShortLinkSetting struct {
	Enabled bool          `yaml:"enabled" json:"enabled"`
	Length  int           `yaml:"length" json:"length"`
	BaseUrl string        `yaml:"baseUrl" json:"baseUrl"`
	Context string        `yaml:"context" json:"context"`
	Expired time.Duration `yaml:"expired" json:"expired"`
}

var Setting = &ShortLinkSetting{
	Enabled: true,
	Length:  5,
	Expired: time.Duration(24*180) * time.Hour, /* 180 DAY */
}

func init() {
	config.RegisterConfig("shortlink", Setting, "ShortLink Config")
}
