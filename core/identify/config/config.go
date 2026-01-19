package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type IdentitySetting struct {
	Enabled bool
}

var Setting = &IdentitySetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Identify Config")
	config.RegisterConfig("identify", Setting, "Identity Settings")
}
