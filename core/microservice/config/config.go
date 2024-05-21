package config

import (
	RegistryConfig "github.com/wjshen/gophrame/core/microservice/registry/config"
)

type MicroserviceSetting struct {
	Enabled  bool
	Registry *RegistryConfig.RegistrySetting `json:"registry"`
}

var Setting *MicroserviceSetting = &MicroserviceSetting{
	Enabled:  false,
	Registry: RegistryConfig.Setting,
}
