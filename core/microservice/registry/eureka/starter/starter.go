package starter

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/microservice/registry/config"
	"github.com/wjshen/gophrame/core/microservice/registry/eureka"
)

func init() {
	if config.Setting.Enabled {
		if client, err := eureka.CreateEurekaDiscoveryClient(); err == nil && client != nil {
			inject.InjectValue("discoveryClient", client)
		}
	}
}
