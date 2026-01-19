package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type DriverSetting struct {
	Enabled               bool          `json:"enabled" yaml:"enabled"`
	SlowThreshold         time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
	MaxIdleConnections    int           `json:"maxIdleConnections" yaml:"maxIdleConnections"`
	ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime" yaml:"connectionMaxIdleTime"`
	MaxOpenConnections    int           `json:"maxOpenConnections" yaml:"maxOpenConnections"`
	ConnectionMaxLifeTime time.Duration `json:"connectionMaxLifeTime" yaml:"connectionMaxLifeTime"`
}

type DatabaseSetting struct {
	// Common Settings
	Driver      string `json:"driver"`
	TablePrefix string `json:"tablePrefix" yaml:"tablePrefix"`
	DriverSetting
	Read *DriverSetting `json:"read,omitempty" yaml:"read"`
}

var Setting *DatabaseSetting = &DatabaseSetting{
	DriverSetting: DriverSetting{
		Enabled: false,

		// 慢SQL时间阈值 = 10s
		SlowThreshold: time.Second * 10,

		// 数据库连接闲置时间 = 30s
		ConnectionMaxIdleTime: time.Second * 30,
		MaxIdleConnections:    10,

		// 数据库连接存在时间 = 180s
		ConnectionMaxLifeTime: time.Second * 180,
		MaxOpenConnections:    128,
	},
}

func init() {
	logger.Debug("Register Database Config")
	config.RegisterConfig("database", Setting, "Database Settings")
}
