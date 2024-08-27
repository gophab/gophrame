package auth

import (
	"strconv"
	"strings"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthService "github.com/gophab/gophrame/default/service/auth"

	"github.com/gin-gonic/gin"
)

type ModuleMController struct {
	controller.ResourceController
	ModuleService *AuthService.ModuleService `inject:"moduleService"`
}

var moduleMController *ModuleMController = &ModuleMController{}

func init() {
	inject.InjectValue("moduleMController", moduleMController)
}

func (m *ModuleMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/modules", Handler: m.GetModules},
		{HttpMethod: "GET", ResourcePath: "/modules/:id", Handler: m.GetSubModules},
		{HttpMethod: "GET", ResourcePath: "/module/:id", Handler: m.GetModule},
		{HttpMethod: "POST", ResourcePath: "/module", Handler: m.CreateModule},
		{HttpMethod: "PUT", ResourcePath: "/module", Handler: m.UpdateModule},
		{HttpMethod: "DELETE", ResourcePath: "/module/:id", Handler: m.DeleteModule},
		{HttpMethod: "GET", ResourcePath: "/tenant/:id/modules", Handler: m.GetTenantModules},
		{HttpMethod: "POST", ResourcePath: "/tenant/:id/modules", Handler: m.AuthorityTenantModules},
		{HttpMethod: "PUT", ResourcePath: "/tenant/:id/modules", Handler: m.SetTenantModules},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id/module/:mid", Handler: m.RemoveTenantModuleAuthority},
	})
}

func (m *ModuleMController) CreateModule(c *gin.Context) {
	var tmp AuthModel.Module
	if err := c.ShouldBind(&tmp); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	tmp.TenantId = "SYSTEM"
	if tmp, err := m.ModuleService.CreateModule(&tmp); err == nil {
		response.Success(c, tmp)
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *ModuleMController) UpdateModule(c *gin.Context) {
	var tmp AuthModel.Module
	if err := c.ShouldBind(&tmp); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	tmp.TenantId = "SYSTEM"
	if tmp, err := m.ModuleService.UpdateModule(&tmp); err == nil {
		response.Success(c, tmp)
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *ModuleMController) DeleteModule(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if err := m.ModuleService.DeleteModule(id); err == nil {
		response.Success(c, "")
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *ModuleMController) GetModule(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, _ := m.ModuleService.GetById(id); res == nil {
		response.NotFound(c, "Not Found")
	} else {
		response.Success(c, res)
	}
}

func (m *ModuleMController) GetModules(c *gin.Context) {
	fid := request.Param(c, "fid").DefaultInt64(-1)
	title := request.Param(c, "title").DefaultString("")
	tree := request.Param(c, "tree").DefaultBool(false)
	pageable := query.GetPageable(c)

	count, result := m.ModuleService.List(fid, title, "SYSTEM", pageable)
	if count != 0 {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10))
		if tree {
			result = m.ModuleService.MakeTree(result)
		}
		response.Success(c, result)
		return
	}
	response.Success(c, []any{})
}

func (m *ModuleMController) GetSubModules(c *gin.Context) {
	fid, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, err := m.ModuleService.GetByFid(fid); err == nil {
		response.Success(c, res)
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

// 待分配的系统菜单以及挂接的按钮
func (c *ModuleMController) GetTenantModules(context *gin.Context) {
	tenantId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	list, err := c.ModuleService.GetTenantModules(tenantId)
	if err != nil {
		response.SystemFail(context, err)
		return
	}
	response.Success(context, list)
}

// 授权企业
func (m *ModuleMController) AuthorityTenantModules(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var moduleIds []int64
	if err := c.ShouldBind(&moduleIds); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	for _, moduleId := range moduleIds {
		// 1. as ID
		if module, err := m.ModuleService.GetById(moduleId); err == nil && module != nil {
			m.ModuleService.AddTenantModule(module, tenantId)
		}
	}

	response.Success(c, nil)
}

// 授权企业
func (m *ModuleMController) SetTenantModules(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var moduleIds []int64
	if err := c.ShouldBind(&moduleIds); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	modules, err := m.ModuleService.GetByIds(moduleIds)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	m.ModuleService.SetTenantModules(tenantId, modules)
	modules, err = m.ModuleService.GetTenantModules(tenantId)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	response.Success(c, modules)
}

// 授权企业
func (m *ModuleMController) RemoveTenantModuleAuthority(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	mid, err := request.Param(c, "mid").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var moduleIds = strings.Split(mid, ",")
	for _, moduleId := range moduleIds {
		// 1. as ID
		if id, err := strconv.Atoi(moduleId); err == nil {
			if module, err := m.ModuleService.GetById(int64(id)); err == nil {
				if module != nil && module.TenantId == tenantId {
					m.ModuleService.DeleteModule(int64(id))
				}
			}
		} else {
			m.ModuleService.DeleteTenantModule(moduleId, tenantId)
		}
	}

	response.Success(c, nil)
}
