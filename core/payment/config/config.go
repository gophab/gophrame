package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type PaymentSetting struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

var Setting *PaymentSetting = &PaymentSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Payment Config")
	config.RegisterConfig("payment", Setting, "Payment Settings")
}
