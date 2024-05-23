package openapi

import (
	"encoding/json"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	SecurityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/wjshen/gophrame/default/domain"
	"github.com/wjshen/gophrame/default/service"

	"github.com/gin-gonic/gin"
)

type TenantOptionOpenController struct {
	controller.ResourceController
	TenantOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var tenantOptionOpenController = &TenantOptionOpenController{}

type AdminTenantOptionOpenController struct {
	controller.ResourceController
	TenantOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var adminTenantOptionOpenController = &AdminTenantOptionOpenController{}

func init() {
	inject.InjectValue("tenantOptionOpenController", tenantOptionOpenController)
	inject.InjectValue("adminTenantOptionOpenController", adminTenantOptionOpenController)
}

func (c *TenantOptionOpenController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant/options", Handler: c.GetTenantOptions},
	})
}

func (c *TenantOptionOpenController) GetTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if tenantOptions, err := c.TenantOptionService.GetTenantOptions(currentTenantId); err == nil {
		result := make(map[string]string)
		for name, option := range tenantOptions.Options {
			if option.Public {
				result[name] = option.Value
			}
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *AdminTenantOptionOpenController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant/options", Handler: c.GetTenantOptions},
		{HttpMethod: "POST", ResourcePath: "/tenant/options", Handler: c.AddTenantOptions},
		{HttpMethod: "PUT", ResourcePath: "/tenant/options", Handler: c.SetTenantOptions},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/option/:key", Handler: c.RemoveTenantOption},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/options", Handler: c.RemoveTenantOptions},
	})
}

func (c *AdminTenantOptionOpenController) GetTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if tenantOptions, err := c.TenantOptionService.GetTenantOptions(currentTenantId); err == nil {
		result := make(map[string]string)
		for name, option := range tenantOptions.Options {
			result[name] = option.Value
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *AdminTenantOptionOpenController) AddTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if body, err := ctx.GetRawData(); err == nil {
		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			if _, err := c.TenantOptionService.AddSysOption(&domain.SysOption{
				TenantId: currentTenantId,
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
		c.GetTenantOptions(ctx)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *AdminTenantOptionOpenController) SetTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if body, err := ctx.GetRawData(); err == nil {
		var tenantOptions = domain.SysOptions{
			TenantId: currentTenantId,
		}

		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			tenantOptions.Options[k] = domain.SysOption{
				TenantId: currentTenantId,
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}
		}

		if _, err := c.TenantOptionService.SetTenantOptions(&tenantOptions); err == nil {
			c.GetTenantOptions(ctx)
		} else {
			response.FailMessage(ctx, 400, err.Error())
		}
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *AdminTenantOptionOpenController) RemoveTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if err := c.TenantOptionService.RemoveAllTenantOptions(currentTenantId); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *AdminTenantOptionOpenController) RemoveTenantOption(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	key := ctx.Param("key")
	if res, err := c.TenantOptionService.RemoveTenantOption(currentTenantId, key); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.Success(ctx, res)
	}
}
