package openapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/operation/service"

	"github.com/gin-gonic/gin"
)

type ModuleOpenController struct {
	controller.ResourceController
	ModuleService *service.ModuleService `inject:"moduleService"`
}

var moduleOpenController *ModuleOpenController = &ModuleOpenController{}

func init() {
	inject.InjectValue("moduleOpenController", moduleOpenController)
}

func (m *ModuleOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/modules", Handler: m.GetCurrentUserModules},
		{HttpMethod: "GET", ResourcePath: "/modules/:id", Handler: m.GetSubModules},
		{HttpMethod: "GET", ResourcePath: "/module/:id", Handler: m.GetModule},
	})
}

// @Summary   获取登录用户信息
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {"lists":""}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/userInfo  [GET]
func (u *ModuleOpenController) GetCurrentUserModules(c *gin.Context) {
	tree := request.Param(c, "tree").DefaultBool(false)
	tenantId := SecurityUtil.GetCurrentTenantId(c)

	if tenantId == "" {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	modules, err := u.ModuleService.GetTenantModules(tenantId)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	if tree {
		response.Success(c, u.ModuleService.MakeTree(modules))
	} else {
		response.Success(c, modules)
	}
}

func (u *ModuleOpenController) GetModule(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	tenantId := SecurityUtil.GetCurrentTenantId(c)
	if tenantId == "" {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	modules, err := u.ModuleService.GetTenantModules(tenantId)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	for _, module := range modules {
		if module.Id == id {
			response.Success(c, module)
			return
		}
	}
	response.NotFound(c, "")
}

func (u *ModuleOpenController) GetSubModules(c *gin.Context) {
	fid, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	tenantId := SecurityUtil.GetCurrentTenantId(c)
	if tenantId == "" {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	modules, err := u.ModuleService.GetTenantModules(tenantId)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	for _, module := range modules {
		if module.Id == fid {
			if res, err := u.ModuleService.GetByFid(fid); err == nil {
				response.Success(c, res)
			} else {
				response.FailMessage(c, 400, err.Error())
			}
			return
		}
	}
	response.NotFound(c, "")
}
