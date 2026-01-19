package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type MysqlSetting struct {
	MysqlDatabase
	Read *MysqlDatabase `json:"read,omitempty" yaml:"read,omitempty"`
}

type MysqlDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
}

var Setting *MysqlSetting = &MysqlSetting{}

func init() {
	logger.Debug("Register Database Config - Mysql")
	config.RegisterConfig("database.mysql", Setting, "Mysql Settings")
}
