package config

import "time"

type MysqlSetting struct {
	Default    *MysqlDatabase `json:"default"`
	EnableRead bool           `json:"enableRead" yaml:"enableRead"`
	Read       *MysqlDatabase `json:"read,omitempty" yaml:"read,omitempty"`
}

type MysqlDatabase struct {
	Host                  string        `json:"host"`
	Port                  int           `json:"port"`
	User                  string        `json:"user"`
	Password              string        `json:"password"`
	Database              string        `json:"database"`
	Charset               string        `json:"charset"`
	MaxIdleConnections    int           `json:"maxIdleConnections" yaml:"maxIdleConnections"`
	ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime" yaml:"connectionMaxIdleTime"`
	MaxOpenConnections    int           `json:"maxOpenConnections" yaml:"maxOpenConnections"`
	ConnectionMaxLifeTime time.Duration `json:"connectionMaxLifeTime" yaml:"connectionMaxLifeTime"`
}

var Setting *MysqlSetting = &MysqlSetting{
	Default: &MysqlDatabase{
		// 数据库连接闲置时间 = 30s
		ConnectionMaxIdleTime: time.Second * 30,
		MaxIdleConnections:    10,

		// 数据库连接存在时间 = 180s
		ConnectionMaxLifeTime: time.Second * 180,
		MaxOpenConnections:    128,
	},
	EnableRead: false,
	Read: &MysqlDatabase{
		// 数据库连接闲置时间 = 30s
		ConnectionMaxIdleTime: time.Second * 30,
		MaxIdleConnections:    10,

		// 数据库连接存在时间 = 180s
		ConnectionMaxLifeTime: time.Second * 180,
		MaxOpenConnections:    128,
	},
}
