package api

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

	"github.com/wjshen/gophrame/default/service"
	"github.com/wjshen/gophrame/default/service/dto"
)

type RoleController struct {
	controller.ResourceController
	RoleService *service.RoleService `inject:"roleController"`
}

var roleController *RoleController = &RoleController{}

func init() {
	inject.InjectValue("roleController", roleController)
}

func (m *RoleController) AfterInitialize() {
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
func (r *RoleController) GetRoles(c *gin.Context) {
	id := com.StrTo(c.Query("id")).String()

	total, err := r.RoleService.Count(&dto.Role{Id: id})
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

	response.Success(c, data)
}

// @Summary   获取角色
// @Tags role
// @Accept
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role/:id  [GET]
func (r *RoleController) GetRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := r.RoleService.Get(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_COUNT_FAIL)
		return
	}

	response.Success(c, result)
}

// @Summary   增加角色
// @Tags role
// @Accept json
// @Produce  json
// @Param   body  body   models.Role   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role  [POST]
func (r *RoleController) AddRole(c *gin.Context) {
	var data dto.RoleCreate
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.MaxSize(data.Name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	role, err := r.RoleService.Add(&dto.Role{
		RoleCreate: data,
	})

	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	response.Success(c, role)
}

// @Summary   更新角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Param   body  body   models.Role   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role/:id  [PUT]
func (r *RoleController) EditRole(c *gin.Context) {
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
		response.FailCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	if exists, err := r.RoleService.ExistByID(data.Id); err != nil {
		response.FailCode(c, errors.ERROR_EXIST_FAIL)
		return
	} else if !exists {
		response.NotFound(c, data.Id)
		return
	}

	if err := r.RoleService.Edit(&data); err == nil {
		r.RoleService.LoadPolicy(data.Id)
		response.Success(c, data)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func (r *RoleController) DeleteRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if exists, err := r.RoleService.ExistByID(id); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_EXIST_FAIL, err.Error())
		return
	} else if !exists {
		response.NotFound(c, id)
		return
	}

	if err := r.RoleService.Delete(id); err == nil {
		response.Success(c, nil)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	}
}
