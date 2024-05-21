package config

import (
	ConsulConfig "github.com/wjshen/gophrame/core/microservice/registry/consul/config"
	DubboConfig "github.com/wjshen/gophrame/core/microservice/registry/dubbo/config"
	EurekaConfig "github.com/wjshen/gophrame/core/microservice/registry/eureka/config"
	NacosConfig "github.com/wjshen/gophrame/core/microservice/registry/nacos/config"

	"github.com/google/uuid"
)

type RegistrySetting struct {
	Enabled            bool
	EnableAutoRegister bool   `json:"enableAutoRegister" yaml:"enabledAutoRegister"`
	ServiceName        string `json:"serviceName" yaml:"serviceName"`
	InstanceId         string `json:"instanceId" yaml:"instanceId"`
	PreferIP           string `json:"perferIp" yaml:"preferIp"`
	Port               int
	Eureka             *EurekaConfig.EurekaSetting
	Consul             *ConsulConfig.ConsulSetting
	Nacos              *NacosConfig.NacosSetting
	Dubbo              *DubboConfig.DubboSetting
}

var Setting = &RegistrySetting{
	Enabled:            false,
	EnableAutoRegister: false,
	InstanceId:         uuid.NewString(),
}
