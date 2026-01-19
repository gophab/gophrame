package openapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"

	"github.com/gin-gonic/gin"
)

var tenantOpenController *TenantOpenController = &TenantOpenController{}
var adminTenantOpenController *AdminTenantOpenController = &AdminTenantOpenController{}

func init() {
	inject.InjectValue("tenantOpenController", tenantOpenController)
	inject.InjectValue("adminTenantOpenController", adminTenantOpenController)
}

type TenantOpenController struct {
	controller.ResourceController
	TenantService *service.TenantService `inject:"tenantService"`
}

type AdminTenantOpenController struct {
	controller.ResourceController
	TenantService *service.TenantService `inject:"tenantService"`
}

// 组织
func (m *TenantOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant", Handler: m.GetTenant},
	})
}

// 1.根据id查询节点
func (a *TenantOpenController) GetTenant(context *gin.Context) {
	tenantId := SecurityUtil.GetCurrentTenantId(context)

	if result, _ := a.TenantService.GetById(tenantId); result != nil {
		response.Success(context, result)
	} else {
		response.NotFound(context, "")
	}
}

// 组织
func (m *AdminTenantOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant", Handler: m.GetTenant},
		{HttpMethod: "PUT", ResourcePath: "/tenant", Handler: m.UpdateTenant},
		{HttpMethod: "PATCH", ResourcePath: "/tenant", Handler: m.PatchTenant},
	})
}

// 1.根据id查询节点
func (a *AdminTenantOpenController) GetTenant(context *gin.Context) {
	tenantId := SecurityUtil.GetCurrentTenantId(context)

	if result, _ := a.TenantService.GetById(tenantId); result != nil {
		response.Success(context, result)
	} else {
		response.NotFound(context, "")
	}
}

// 修改
func (a *AdminTenantOpenController) UpdateTenant(c *gin.Context) {
	var data domain.Tenant
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.TenantService.Update(&data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 修改
func (a *AdminTenantOpenController) PatchTenant(c *gin.Context) {
	var data = make(map[string]any)
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var availableFields = []string{"name", "description", "telephone", "address", "logo", "status"}
	var tenant = make(map[string]any)
	for _, k := range availableFields {
		if v, b := data[k]; b && v != nil {
			switch t := v.(type) {
			case string:
				if t != "" {
					tenant[k] = v
				}
			default:
				tenant[k] = v
			}
		}
	}

	var id = data["id"].(string)
	if result, err := a.TenantService.PatchAll(id, tenant); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}
