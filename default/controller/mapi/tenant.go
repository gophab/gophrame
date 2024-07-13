package mapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/service"

	"github.com/gin-gonic/gin"
)

var tenantMController *TenantMController = &TenantMController{}

func init() {
	inject.InjectValue("tenantMController", tenantMController)
}

type TenantMController struct {
	controller.ResourceController
	TenantService *service.TenantService `inject:"tenantService"`
}

// 组织
func (m *TenantMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant/:id", Handler: m.GetTenant},
		{HttpMethod: "GET", ResourcePath: "/tenants", Handler: m.GetTenants},
		{HttpMethod: "POST", ResourcePath: "/tenant", Handler: m.CreateTenant},
		{HttpMethod: "PUT", ResourcePath: "/tenant", Handler: m.UpdateTenant},
		{HttpMethod: "PATCH", ResourcePath: "/tenant/:id", Handler: m.PatchTenant},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id", Handler: m.DeleteTenant},
	})
}

// 1.根据id查询节点
func (a *TenantMController) GetTenant(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, _ := a.TenantService.GetById(id); result != nil {
		response.Success(c, result)
	} else {
		response.NotFound(c, "")
	}
}

// 1.根据参数查询租户
func (a *TenantMController) GetTenants(c *gin.Context) {
	search := request.Param(c, "search").DefaultString("")
	id := request.Param(c, "id").DefaultString("")
	name := request.Param(c, "name").DefaultString("")
	licenseId := request.Param(c, "licenseId").DefaultString("")

	pageable := query.GetPageable(c)

	conds := make(map[string]interface{})
	if search != "" {
		conds["search"] = search
	}
	if id != "" {
		conds["id"] = id
	}
	if name != "" {
		conds["name"] = name
	}
	if licenseId != "" {
		conds["licenseId"] = licenseId
	}

	if total, result := a.TenantService.Find(conds, pageable); result != nil {
		response.Page(c, total, result)
	} else {
		response.NotFound(c, "")
	}
}

// 创建
func (a *TenantMController) CreateTenant(c *gin.Context) {
	var data domain.Tenant
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if data.Id != "" {
		if result, err := a.TenantService.Update(&data); err == nil {
			response.Success(c, result)
		} else {
			response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
		}
	} else {
		if result, err := a.TenantService.Create(&data); err == nil {
			response.Success(c, result)
		} else {
			response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
		}
	}
}

// 修改
func (a *TenantMController) UpdateTenant(c *gin.Context) {
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
func (a *TenantMController) PatchTenant(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var data = make(map[string]interface{})
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.TenantService.PatchAll(id, data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 删除
func (a *TenantMController) DeleteTenant(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.TenantService.DeleteById(id); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}
