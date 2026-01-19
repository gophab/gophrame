package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type MongoSetting struct {
	Enabled                bool          `json:"enabled"`
	Host                   string        `json:"host"`
	Port                   int           `json:"port"`
	User                   string        `json:"user"`
	Password               string        `json:"password"`
	AdminDatabase          string        `json:"adminDatabase"`
	Database               string        `json:"database"`
	MaxPoolSize            uint64        `json:"maxPoolSize"`            /* 最大连接数 */
	MinPoolSize            uint64        `json:"minPoolSize"`            /* 最小连接数 */
	MaxConnIdleTime        time.Duration `json:"maxConnIdleTime"`        /* 连接最大空闲时间 */
	ConnectTimeout         time.Duration `json:"connectTimeout"`         /* 连接超时时间 */
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout"` /* 服务器选择超时时间 */
}

var Setting = &MongoSetting{
	Enabled:                false,
	MaxPoolSize:            20,
	MinPoolSize:            5,
	MaxConnIdleTime:        5 * time.Minute,
	ConnectTimeout:         10 * time.Second,
	ServerSelectionTimeout: 10 * time.Second,
}

func init() {
	logger.Debug("Register Mongo Config")
	config.RegisterConfig("mongo", Setting, "Mongo Settings")
}
