package mapi

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/collection"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"
)

type RoleMController struct {
	controller.ResourceController
	RoleService     *service.RoleService     `inject:"roleService"`
	RoleUserService *service.RoleUserService `inject:"roleUserService"`
	TenantService   *service.TenantService   `inject:"tenantService"`
}

var roleMController *RoleMController = &RoleMController{}

func init() {
	inject.InjectValue("roleMController", roleMController)
}

func (m *RoleMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/roles", Handler: m.GetRoles},
		{HttpMethod: "GET", ResourcePath: "/role/:id", Handler: m.GetRole},
		{HttpMethod: "POST", ResourcePath: "/role", Handler: m.CreateRole},
		{HttpMethod: "POST", ResourcePath: "/role/:id/copy", Handler: m.CopyRole},
		{HttpMethod: "PUT", ResourcePath: "/role", Handler: m.UpdateRole},
		{HttpMethod: "DELETE", ResourcePath: "/role/:id", Handler: m.DeleteRole},
		{HttpMethod: "GET", ResourcePath: "/role/:id/users", Handler: m.GetRoleUsers},
		{HttpMethod: "POST", ResourcePath: "/role/:id/users", Handler: m.AddRoleUsers},
		{HttpMethod: "DELETE", ResourcePath: "/role/:id/users", Handler: m.DeleteRoleUsers},
	})
}

// @Summary   获取所有角色
// @Tags role
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles  [GET]
func (r *RoleMController) GetRoles(c *gin.Context) {
	pageable := query.GetPageable(c)

	count, roles, err := r.RoleService.Find(map[string]any{
		"del_flag":  false,
		"tenant_id": "SYSTEM",
	}, pageable)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(count, 10))
	response.Success(c, roles)
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

	result, err := r.RoleService.GetById(id)
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
func (r *RoleMController) CreateRole(c *gin.Context) {
	var data domain.RoleInfo
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

	role := &domain.Role{
		RoleInfo: data,
	}
	role.TenantId = "SYSTEM"

	role, err := r.RoleService.CreateRole(role)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	response.Success(c, role)
}

// @Summary   获取角色
// @Tags role
// @Accept
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/role/:id  [GET]
func (r *RoleMController) CopyRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := r.RoleService.GetById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_COUNT_FAIL)
		return
	}

	copied := &domain.Role{RoleInfo: result.RoleInfo}
	copied.Name = copied.Name + "_" + time.Now().Format("20060102150405")
	copied.TenantId = "SYSTEM"

	role, err := r.RoleService.CreateRole(copied)
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
func (r *RoleMController) UpdateRole(c *gin.Context) {
	var data domain.Role
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

	exists, err := r.RoleService.ExistById(data.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	result, err := r.RoleService.UpdateRole(&data)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
		return
	}

	response.Success(c, result)
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

	exists, err := r.RoleService.ExistById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	err = r.RoleService.DeleteById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.Success(c, nil)
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func (u *RoleMController) GetRoleUsers(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := u.RoleService.ExistById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	search := request.Param(c, "search").DefaultString("")
	tenantId := request.Param(c, "tenantId").DefaultString("")
	pageable := query.GetPageable(c)

	count, list := u.RoleUserService.ListUsers(id, search, tenantId, pageable)

	tenantIds := collection.MapToSet[string](list, func(i any) string {
		return i.(*domain.User).TenantId
	})

	var tenants = make(map[string]*domain.Tenant)
	if list, err := u.TenantService.GetByIds(tenantIds.AsList()); err == nil {
		for _, item := range list {
			tenants[item.Id] = item
		}
	}
	tenants["SYSTEM"] = &domain.Tenant{
		Id:   "SYSTEM",
		Name: "平台",
	}
	for _, v := range list {
		v.Password = ""
		v.Tenant = tenants[v.TenantId]
	}

	response.Page(c, count, list)
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func (r *RoleMController) AddRoleUsers(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var userIds []string
	if err := c.ShouldBind(&userIds); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := r.RoleService.ExistById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	results, err := r.RoleUserService.AddRoleUserIds(id, userIds)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.Success(c, results)
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func (r *RoleMController) DeleteRoleUsers(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	users, err := request.Param(c, "users").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var userIds = strings.Split(users, ",")

	exists, err := r.RoleService.ExistById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	err = r.RoleUserService.DeleteRoleUserIds(id, userIds)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.Success(c, nil)
}
