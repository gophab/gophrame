package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type PostgresSetting struct {
	PostgresDatabase
	Read *PostgresDatabase `json:"read,omitempty" yaml:"read,omitempty"`
}

type PostgresDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	TimeZone string `json:"timeZone" yaml:"timeZone"`
}

var Setting *PostgresSetting = &PostgresSetting{}

func init() {
	logger.Debug("Register Database Config - Postgres")
	config.RegisterConfig("database.postgres", Setting, "Postgres Settings")
}
