package config

import (
	"time"

	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"
)

type RedisSetting struct {
	Enabled                  bool          `json:"enabled"`
	Host                     string        `json:"host"`
	Port                     int           `json:"port"`
	Auth                     string        `json:"auth"`
	MaxIdle                  int           `json:"maxIdle" yaml:"maxIdle"`
	MaxActive                int           `json:"maxActive" yaml:"maxActive"`
	IdleTimout               time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
	Database                 int           `json:"database"`
	ConnectionFailRetryTimes int           `json:"connectionFailRetryTimes" yaml:"connectionFailRetryTimes"`
	ReConnectInterval        time.Duration `json:"reConnectInterval" yaml:"reConnectInterval"`
}

var Setting *RedisSetting = &RedisSetting{
	Enabled:                  false,
	Port:                     6379,
	Database:                 1,
	MaxIdle:                  10,
	MaxActive:                1000,
	IdleTimout:               time.Second * 60,
	ConnectionFailRetryTimes: 3,
	ReConnectInterval:        time.Second * 5,
}

func init() {
	logger.Debug("Register Redis Config")
	config.RegisterConfig("redis", Setting, "Redis Settings")
}
