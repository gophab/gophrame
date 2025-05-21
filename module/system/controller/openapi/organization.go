package openapi

import (
	"strconv"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"

	"github.com/gin-gonic/gin"
)

var organizationOpenController *OrganizationOpenController = &OrganizationOpenController{}
var adminOrganizationOpenController *AdminOrganizationOpenController = &AdminOrganizationOpenController{}

func init() {
	inject.InjectValue("organizationOpenController", organizationOpenController)
	inject.InjectValue("adminOrganizationOpenController", adminOrganizationOpenController)
}

type OrganizationOpenController struct {
	controller.ResourceController
	OrganizationService *service.OrganizationService `inject:"organizationService"`
}

type AdminOrganizationOpenController struct {
	controller.ResourceController
	OrganizationService *service.OrganizationService `inject:"organizationService"`
}

// 组织
func (m *OrganizationOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/organizations", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/organizations/:id", Handler: m.GetSubList},
		{HttpMethod: "GET", ResourcePath: "/organization/:id", Handler: m.GetOrganization},
	})
}

// 1.根据id查询节点
func (a *OrganizationOpenController) GetOrganization(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if result, _ := a.OrganizationService.GetById(id); result != nil {
		response.Success(context, result)
	} else {
		response.NotFound(context, "")
	}
}

// 1.省份城市列表
func (a *OrganizationOpenController) GetList(context *gin.Context) {
	fid := request.Param(context, "fid").DefaultString("")
	name := request.Param(context, "name").DefaultString("")
	pageable := query.GetPageable(context)

	if counts, lists := a.OrganizationService.List(fid, name, pageable); counts > 0 {
		context.Header("X-Total-Count", strconv.FormatInt(counts, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 1.根据fid查询子节点列表
func (a *OrganizationOpenController) GetSubList(context *gin.Context) {
	fid, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if subList := a.OrganizationService.GetSubList(fid); len(subList) > 0 {
		response.Success(context, subList)
	} else {
		response.Success(context, []any{})
	}
}

// 组织
func (m *AdminOrganizationOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/organizations", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/organizations/:id", Handler: m.GetSubList},
		{HttpMethod: "GET", ResourcePath: "/organization/:id", Handler: m.GetOrganization},
		{HttpMethod: "POST", ResourcePath: "/organization", Handler: m.CreateOrganization},
		{HttpMethod: "PUT", ResourcePath: "/organization", Handler: m.UpdateOrganization},
		{HttpMethod: "DELETE", ResourcePath: "/organization/:id", Handler: m.DeleteOrganization},
	})
}

// 1.根据id查询节点
func (a *AdminOrganizationOpenController) GetOrganization(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if result, _ := a.OrganizationService.GetById(id); result != nil {
		response.Success(context, result)
	} else {
		response.NotFound(context, "")
	}
}

// 1.省份城市列表
func (a *AdminOrganizationOpenController) GetList(context *gin.Context) {
	fid := request.Param(context, "fid").DefaultString("")
	name := request.Param(context, "name").DefaultString("")
	pageable := query.GetPageable(context)

	if counts, lists := a.OrganizationService.List(fid, name, pageable); counts > 0 {
		context.Header("X-Total-Count", strconv.FormatInt(counts, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 1.根据fid查询子节点列表
func (a *AdminOrganizationOpenController) GetSubList(context *gin.Context) {
	fid, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if subList := a.OrganizationService.GetSubList(fid); len(subList) > 0 {
		response.Success(context, subList)
	} else {
		response.Success(context, []any{})
	}
}

// 新增
func (a *AdminOrganizationOpenController) CreateOrganization(c *gin.Context) {
	var data domain.Organization
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.OrganizationService.CreateOrganization(&data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// 修改
func (a *AdminOrganizationOpenController) UpdateOrganization(c *gin.Context) {
	var data domain.Organization
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.OrganizationService.UpdateOrganization(&data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 删除
func (a *AdminOrganizationOpenController) DeleteOrganization(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if a.OrganizationService.HasSubNode(id) > 0 {
		response.FailMessage(c, errors.ERROR_DELETE_FAIL, "该节点下有子节点,禁止删除")
		return
	}

	if b, err := a.OrganizationService.DeleteOrganization(id); b {
		response.Success(c, "")
	} else {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	}
}
