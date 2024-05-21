package openapi

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/service"

	"github.com/gin-gonic/gin"
)

type PublicSystemOptionOpenController struct {
	controller.ResourceController
	SystemOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var publicSystemOptionOpenController = &PublicSystemOptionOpenController{}

func init() {
	inject.InjectValue("publicSystemOptionOpenController", publicSystemOptionOpenController)

}

func (c *PublicSystemOptionOpenController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/system/options", Handler: c.GetSystemOptions},
	})
}

// GET /system/options
func (c *PublicSystemOptionOpenController) GetSystemOptions(ctx *gin.Context) {
	if systemOptions, err := c.SystemOptionService.GetSystemOptions(); err == nil {
		result := make(map[string]string)
		for name, option := range systemOptions.Options {
			if option.Public {
				result[name] = option.Value
			}
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}
