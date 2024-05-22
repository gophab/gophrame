package starter

import (
	_ "github.com/wjshen/gophrame/core/microservice/registry/consul/starter"
	_ "github.com/wjshen/gophrame/core/microservice/registry/dubbo/starter"
	_ "github.com/wjshen/gophrame/core/microservice/registry/eureka/starter"
	_ "github.com/wjshen/gophrame/core/microservice/registry/nacos/starter"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/microservice/registry"
	"github.com/wjshen/gophrame/core/microservice/registry/config"
	"github.com/wjshen/gophrame/core/starter"
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
