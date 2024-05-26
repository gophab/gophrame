package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"

	CodeConfig "github.com/gophab/gophrame/core/code/config"
)

var Setting = &CodeConfig.CodeStoreSetting{
	Enabled:         true,
	RequestInterval: time.Minute,
	ExpireIn:        time.Hour * 3 * 24,
	Redis:           nil,
}

func init() {
	logger.Debug("Register Email Code Store Config")
	config.RegisterConfig("email.store", Setting, "Email Code Store Settings")
}
