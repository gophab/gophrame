package openapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/json"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"

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
		{HttpMethod: "PATCH", ResourcePath: "/tenant/options", Handler: c.UpdateTenantOptions},
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
		var data = make(map[string]interface{})
		_ = json.Json(string(body), &data)
		for k, v := range data {
			if _, err := c.TenantOptionService.AddSysOption(&domain.SysOption{
				TenantId: currentTenantId,
				Option: domain.Option{
					Name:      k,
					Value:     json.String(v),
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
			Options:  make(map[string]*domain.SysOption),
		}

		var data = make(map[string]interface{})
		_ = json.Json(string(body), &data)
		for k, v := range data {
			tenantOptions.Options[k] = &domain.SysOption{
				TenantId: currentTenantId,
				Option: domain.Option{
					Name:      k,
					Value:     json.String(v),
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

func (c *AdminTenantOptionOpenController) UpdateTenantOptions(ctx *gin.Context) {
	currentTenantId := SecurityUtil.GetCurrentTenantId(ctx)
	if currentTenantId == "" {
		response.FailMessage(ctx, 403, "")
		return
	}

	if body, err := ctx.GetRawData(); err == nil {
		var tenantOptions = domain.SysOptions{
			TenantId: currentTenantId,
			Options:  make(map[string]*domain.SysOption),
		}

		var data = make(map[string]interface{})
		_ = json.Json(string(body), &data)
		for k, v := range data {
			tenantOptions.Options[k] = &domain.SysOption{
				TenantId: currentTenantId,
				Option: domain.Option{
					Name:      k,
					Value:     json.String(v),
					ValueType: "STRING",
				},
			}
		}

		if _, err := c.TenantOptionService.UpdateTenantOptions(&tenantOptions); err == nil {
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
