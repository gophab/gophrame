package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type ServerSetting struct {
	Enabled          bool          `json:"enabled" yaml:"enabled"`
	BindAddr         string        `json:"bindAddr" yaml:"bindAddr"`
	Port             int           `json:"port"`
	ReadTimeout      time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout     time.Duration `json:"wirteTimeout" yaml:"writeTimeout"`
	AllowCrossDomain bool          `json:"allowCrossDomain" yaml:"allowCrossDomain"`
}

var Setting *ServerSetting = &ServerSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Server Config")
	config.RegisterConfig("server", Setting, "Server Settings")
}
