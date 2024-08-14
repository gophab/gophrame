package openapi

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
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/service"
)

type AdminRoleOpenController struct {
	controller.ResourceController
	RoleService     *service.RoleService     `inject:"roleService"`
	RoleUserService *service.RoleUserService `inject:"roleUserService"`
	UserService     *service.UserService     `inject:"userService"`
	TenantService   *service.TenantService   `inject:"tenantService"`
}

var adminRoleOpenController *AdminRoleOpenController = &AdminRoleOpenController{}

func init() {
	inject.InjectValue("roleOpenController", adminRoleOpenController)
}

func (m *AdminRoleOpenController) AfterInitialize() {
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
func (r *AdminRoleOpenController) GetRoles(c *gin.Context) {
	pageable := query.GetPageable(c)

	count, roles, err := r.RoleService.FindAvailable(map[string]interface{}{
		"del_flag":  false,
		"tenant_id": SecurityUtil.GetCurrentTenantId(c),
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
func (r *AdminRoleOpenController) GetRole(c *gin.Context) {
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

	if result == nil {
		response.NotFound(c, "Not Found")
		return
	}

	if result.TenantId != SecurityUtil.GetCurrentTenantId(c) && result.Scope != "PUBLIC" {
		response.NotAllowed(c, "Not Allowed")
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
func (r *AdminRoleOpenController) CreateRole(c *gin.Context) {
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
	role.TenantId = SecurityUtil.GetCurrentTenantId(c)

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
func (r *AdminRoleOpenController) CopyRole(c *gin.Context) {
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
	copied.TenantId = SecurityUtil.GetCurrentTenantId(c)

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
func (r *AdminRoleOpenController) UpdateRole(c *gin.Context) {
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

	exist, err := r.RoleService.GetById(data.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if exist == nil {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	if exist.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
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
func (r *AdminRoleOpenController) DeleteRole(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exist, err := r.RoleService.GetById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if exist == nil {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	if exist.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
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
func (u *AdminRoleOpenController) GetRoleUsers(c *gin.Context) {
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
	tenantId := SecurityUtil.GetCurrentTenantId(c)
	pageable := query.GetPageable(c)

	count, list := u.RoleUserService.ListUsers(id, search, tenantId, pageable)
	for _, v := range list {
		v.Password = ""
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
func (r *AdminRoleOpenController) AddRoleUsers(c *gin.Context) {
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

	var filterUserIds = make([]string, 0)
	for _, userId := range userIds {
		if user, err := r.UserService.GetById(userId); err == nil {
			if user.TenantId == SecurityUtil.GetCurrentTenantId(c) {
				filterUserIds = append(filterUserIds, userId)
			}
		}
	}

	if len(filterUserIds) > 0 {
		results, err := r.RoleUserService.AddRoleUserIds(id, filterUserIds)
		if err != nil {
			response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
			return
		}
		response.Success(c, results)
	}

	response.Success(c, []*domain.RoleUser{})
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func (r *AdminRoleOpenController) DeleteRoleUsers(c *gin.Context) {
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

	var filterUserIds = make([]string, 0)
	for _, userId := range userIds {
		if user, err := r.UserService.GetById(userId); err == nil {
			if user.TenantId == SecurityUtil.GetCurrentTenantId(c) {
				filterUserIds = append(filterUserIds, userId)
			}
		}
	}

	if len(filterUserIds) > 0 {
		err = r.RoleUserService.DeleteRoleUserIds(id, filterUserIds)
		if err != nil {
			response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
			return
		}
	}

	response.Success(c, nil)
}
