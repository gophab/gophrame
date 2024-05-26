package starter

import (
	_ "github.com/gophab/gophrame/core/microservice/registry/consul/starter"
	_ "github.com/gophab/gophrame/core/microservice/registry/dubbo/starter"
	_ "github.com/gophab/gophrame/core/microservice/registry/eureka/starter"
	_ "github.com/gophab/gophrame/core/microservice/registry/nacos/starter"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/microservice/registry"
	"github.com/gophab/gophrame/core/microservice/registry/config"
	"github.com/gophab/gophrame/core/starter"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	if config.Setting.Enabled {
		// 启动RegistryClient
		registryClient := registry.NewRegistryClient()
		inject.InjectValue("registryClient", registryClient)

		registryClient.Init()
	}
}
