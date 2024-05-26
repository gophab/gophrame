package config

import (
	"time"

	CodeConfig "github.com/gophab/gophrame/core/code/config"
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

var Setting = &CodeConfig.CodeStoreSetting{
	Enabled:         true,
	RequestInterval: time.Minute,
	ExpireIn:        time.Minute * 5,
	Redis:           nil,
}

func init() {
	logger.Debug("Register SMS Code Store Config")
	config.RegisterConfig("sms.store", Setting, "SMS Code Store Settings")
}
