package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type EmailSenderSetting struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	AuthUser string `json:"authUser" yaml:"authUser"`
	AuthPass string `json:"authPass" yaml:"authPass"`
	From     string `json:"from" yaml:"from"`
}

type EmailSetting struct {
	Enabled bool                `json:"enabled" yaml:"enabled"`
	Sender  *EmailSenderSetting `json:"sender" yaml:"sender"`
}

var Setting *EmailSetting = &EmailSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Email Config")
	config.RegisterConfig("email", Setting, "Email Settings")
}
