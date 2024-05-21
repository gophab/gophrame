package config

import (
	"time"

	MysqlConfig "github.com/wjshen/gophrame/core/database/mysql/config"
)

type DatabaseSetting struct {
	// Common Settings
	Enabled               bool          `json:"enabled"`
	Driver                string        `json:"driver"`
	TablePrefix           string        `json:"tablePrefix" yaml:"tablePrefix"`
	SlowThreshold         time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
	MaxIdleConnections    int           `json:"maxIdleConnections" yaml:"maxIdleConnections"`
	ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime" yaml:"connectionMaxIdleTime"`
	MaxOpenConnections    int           `json:"maxOpenConnections" yaml:"maxOpenConnections"`
	ConnectionMaxLifeTime time.Duration `json:"connectionMaxLifeTime" yaml:"connectionMaxLifeTime"`

	// Driver Settings
	Mysql *MysqlConfig.MysqlSetting `json:"mysql"`
}

var Setting *DatabaseSetting = &DatabaseSetting{
	Enabled: false,

	// 慢SQL时间阈值 = 10s
	SlowThreshold: time.Second * 10,

	// 数据库连接闲置时间 = 30s
	ConnectionMaxIdleTime: time.Second * 30,
	MaxIdleConnections:    10,

	// 数据库连接存在时间 = 180s
	ConnectionMaxLifeTime: time.Second * 180,
	MaxOpenConnections:    128,

	Mysql: MysqlConfig.Setting,
}
