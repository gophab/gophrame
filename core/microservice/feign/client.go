package feign

import (
	"net/url"
	"time"

	"github.com/wjshen/gophrame/core/feign"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/microservice/registry"
	"github.com/wjshen/gophrame/core/microservice/registry/config"

	"github.com/patrickmn/go-cache"
)

type RegistryFeignClientInterceptor struct {
	RegistryClient registry.Client
	Cache          *cache.Cache
}

func (in *RegistryFeignClientInterceptor) Do(chain *feign.FeignClientInterceptorChain, method string, urlPath string, urlValues url.Values, bodyValue interface{}, options ...*feign.RequestOptions) *feign.FeignClient {
	if in.RegistryClient != nil {
		url, err := url.Parse(urlPath)
		if err == nil && url.Host != "" {
			if v, b := in.getCache().Get(url.Host); b {
				urlPrefix := v.(string)
				if urlPrefix != "-" {
					// - 标识非服务名，忽略
					urlPath = urlPrefix + url.Path
				}
			} else {
				// 假设Host为ServiceName
				urlPrefix, _ := in.RegistryClient.GetServiceEntry(url.Host)
				if urlPrefix != "" {
					urlPath = urlPrefix + url.Path
					in.getCache().SetDefault(url.Host, urlPrefix)
				} else {
					in.getCache().SetDefault(url.Host, "-")
				}
			}
		}
	}
	return chain.Next(method, urlPath, urlValues, bodyValue, options...)
}

func (in *RegistryFeignClientInterceptor) getCache() *cache.Cache {
	if in.Cache == nil {
		in.Cache = cache.New(time.Minute*5, time.Minute)
	}
	return in.Cache
}

func init() {
	if config.Setting.Enabled {
		var interceptor = &RegistryFeignClientInterceptor{}
		inject.InjectValue("registryFeignInterceptor", interceptor)
		feign.RegisterGlobalFeignClientInterceptor(interceptor)
	}
}
