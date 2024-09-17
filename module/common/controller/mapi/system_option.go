package mapi

import (
	"encoding/json"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"

	"github.com/gin-gonic/gin"
)

type SystemOptionMController struct {
	controller.ResourceController
	SystemOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var systemOptionMController = &SystemOptionMController{}

func init() {
	inject.InjectValue("systemOptionMController", systemOptionMController)
}

func (c *SystemOptionMController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/system/options", Handler: c.GetSystemOptions},
		{HttpMethod: "POST", ResourcePath: "/system/options", Handler: c.AddSystemOptions},
		{HttpMethod: "PUT", ResourcePath: "/system/options", Handler: c.SetSystemOptions},
		{HttpMethod: "DELETE", ResourcePath: "/system/option/:key", Handler: c.RemoveSystemOption},
		{HttpMethod: "DELETE", ResourcePath: "/system/options", Handler: c.RemoveSystemOptions},
	})
}

func (c *SystemOptionMController) GetSystemOptions(ctx *gin.Context) {
	if systemOptions, err := c.SystemOptionService.GetSystemOptions(); err == nil {
		result := make(map[string]string)
		for name, option := range systemOptions.Options {
			result[name] = option.Value
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *SystemOptionMController) AddSystemOptions(ctx *gin.Context) {
	if body, err := ctx.GetRawData(); err == nil {
		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			if _, err := c.SystemOptionService.AddSysOption(&domain.SysOption{
				TenantId: "SYSTEM",
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}); err != nil {
				response.FailMessage(ctx, 400, err.Error())
				return
			}
		}
		c.GetSystemOptions(ctx)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *SystemOptionMController) SetSystemOptions(ctx *gin.Context) {
	if body, err := ctx.GetRawData(); err == nil {
		var systemOptions = domain.SysOptions{
			TenantId: "SYSTEM",
		}

		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			systemOptions.Options[k] = domain.SysOption{
				TenantId: "SYSTEM",
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}
		}

		if _, err := c.SystemOptionService.SetTenantOptions(&systemOptions); err == nil {
			c.GetSystemOptions(ctx)
		} else {
			response.FailMessage(ctx, 400, err.Error())
		}
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *SystemOptionMController) RemoveSystemOptions(ctx *gin.Context) {
	if err := c.SystemOptionService.RemoveAllTenantOptions("SYSTEM"); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *SystemOptionMController) RemoveSystemOption(ctx *gin.Context) {
	key := ctx.Param("key")
	if res, err := c.SystemOptionService.RemoveTenantOption("SYSTEM", key); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.Success(ctx, res)
	}
}
