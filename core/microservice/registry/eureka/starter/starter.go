package starter

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/microservice/registry/config"
	"github.com/gophab/gophrame/core/microservice/registry/eureka"
)

func init() {
	if config.Setting.Enabled {
		if client, err := eureka.CreateEurekaDiscoveryClient(); err == nil && client != nil {
			inject.InjectValue("discoveryClient", client)
		}
	}
}
