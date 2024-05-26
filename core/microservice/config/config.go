package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"

	RegistryConfig "github.com/gophab/gophrame/core/microservice/registry/config"
)

type MicroserviceSetting struct {
	Enabled  bool
	Registry *RegistryConfig.RegistrySetting `json:"registry"`
}

var Setting *MicroserviceSetting = &MicroserviceSetting{
	Enabled:  false,
	Registry: RegistryConfig.Setting,
}

func init() {
	logger.Debug("Register Microservice Config")
	config.RegisterConfig("microservice", Setting, "Microservice Settings")
}
