package mapi

import (
	"encoding/json"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/errors"
	"github.com/wjshen/gophrame/service"

	"github.com/gin-gonic/gin"
)

type TenantOptionMController struct {
	controller.ResourceController
	TenantOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var tenantOptionMController = &TenantOptionMController{}

func init() {
	inject.InjectValue("tenantOptionMController", tenantOptionMController)
}

func (c *TenantOptionMController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant/:id/options", Handler: c.GetTenantOptions},
		{HttpMethod: "POST", ResourcePath: "/tenant/:id/options", Handler: c.AddTenantOptions},
		{HttpMethod: "PUT", ResourcePath: "/tenant/:id/options", Handler: c.SetTenantOptions},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id/option/:key", Handler: c.RemoveTenantOption},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id/options", Handler: c.RemoveTenantOptions},
	})
}

func (c *TenantOptionMController) GetTenantOptions(ctx *gin.Context) {
	tenantId, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.Fail(ctx, 400, errors.GetErrorMessage(errors.INVALID_PARAMS))
		return
	}

	if tenantOptions, err := c.TenantOptionService.GetTenantOptions(tenantId); err == nil {
		result := make(map[string]string)
		for name, option := range tenantOptions.Options {
			result[name] = option.Value
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *TenantOptionMController) AddTenantOptions(ctx *gin.Context) {
	tenantId, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.Fail(ctx, 400, errors.GetErrorMessage(errors.INVALID_PARAMS))
		return
	}

	if body, err := ctx.GetRawData(); err == nil {
		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			if _, err := c.TenantOptionService.AddSysOption(&domain.SysOption{
				TenantId: tenantId,
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

func (c *TenantOptionMController) SetTenantOptions(ctx *gin.Context) {
	tenantId, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.Fail(ctx, 400, errors.GetErrorMessage(errors.INVALID_PARAMS))
		return
	}

	if body, err := ctx.GetRawData(); err == nil {
		var tenantOptions = domain.SysOptions{
			TenantId: tenantId,
		}

		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			tenantOptions.Options[k] = domain.SysOption{
				TenantId: tenantId,
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

func (c *TenantOptionMController) RemoveTenantOptions(ctx *gin.Context) {
	tenantId, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.Fail(ctx, 400, errors.GetErrorMessage(errors.INVALID_PARAMS))
		return
	}

	if err := c.TenantOptionService.RemoveAllTenantOptions(tenantId); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.OK(ctx, nil)
	}
}

func (c *TenantOptionMController) RemoveTenantOption(ctx *gin.Context) {
	tenantId, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.Fail(ctx, 400, errors.GetErrorMessage(errors.INVALID_PARAMS))
		return
	}

	key := ctx.Param("key")
	if res, err := c.TenantOptionService.RemoveTenantOption(tenantId, key); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.OK(ctx, res)
	}
}
