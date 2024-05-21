package mapi

import (
	"github.com/unknwon/com"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"
	"github.com/wjshen/gophrame/service"
	"github.com/wjshen/gophrame/service/dto"
)

type RoleMController struct {
	controller.ResourceController
	RoleService *service.RoleService `inject:"roleMController"`
}

var roleMController *RoleMController = &RoleMController{}

func init() {
	inject.InjectValue("roleMController", roleMController)
}

func (m *RoleMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/roles", Handler: m.GetRoles},
		{HttpMethod: "GET", ResourcePath: "/role/:id", Handler: m.GetRole},
		{HttpMethod: "POST", ResourcePath: "/role", Handler: m.AddRole},
		{HttpMethod: "PUT", ResourcePath: "/role", Handler: m.EditRole},
		{HttpMethod: "DELETE", ResourcePath: "/role/:id", Handler: m.DeleteRole},
	})
}

// @Summary   获取所有角色
// @Tags role
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles  [GET]
func (r *RoleMController) GetRoles(c *gin.Context) {
	id := com.StrTo(c.Query("id")).String()

	total, err := service.GetRoleService().Count(&dto.Role{Id: id})
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_COUNT_FAIL)
		return
	}

	roles, err := r.RoleService.GetAll(&dto.Role{Id: id}, query.GetPageable(c))
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = roles
	data["total"] = total

	response.OK(c, data)
}

// @Summary   获取角色
// @Tags role
// @Accept
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role/:id  [GET]
func (r *RoleMController) GetRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := service.GetRoleService().Get(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_COUNT_FAIL)
		return
	}

	response.OK(c, result)
}

// @Summary   增加角色
// @Tags role
// @Accept json
// @Produce  json
// @Param   body  body   models.Role   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role  [POST]
func (r *RoleMController) AddRole(c *gin.Context) {
	var data dto.RoleCreate
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.MaxSize(data.Name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	role, err := r.RoleService.Add(&dto.Role{
		RoleCreate: data,
	})

	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	response.OK(c, role)
}

// @Summary   更新角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Param   body  body   models.Role   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role/:id  [PUT]
func (r *RoleMController) EditRole(c *gin.Context) {
	var data dto.Role
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.Required(data.Id, "id").Message("Id不能为空")
	valid.MaxSize(data.Name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.FailCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	exists, err := service.GetRoleService().ExistByID(data.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	err = service.GetRoleService().Edit(&data)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EDIT_FAIL)
		return
	}

	err = service.GetRoleService().LoadPolicy(data.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EDIT_FAIL)
		return
	}

	response.OK(c, nil)
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func (r *RoleMController) DeleteRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := service.GetRoleService().ExistByID(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	err = service.GetRoleService().Delete(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.OK(c, nil)
}
