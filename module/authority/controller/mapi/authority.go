package mapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/authority/service"

	OperationModel "github.com/gophab/gophrame/module/operation/domain"

	"github.com/gin-gonic/gin"
)

type AuthorityMController struct {
	controller.ResourceController
	AuthorityService *service.AuthorityService `inject:"authorityService"`
}

var authorityMController *AuthorityMController = &AuthorityMController{}

func init() {
	inject.InjectValue("authorityMController", authorityMController)
}

func (m *AuthorityMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/operations", Handler: m.GetSystemOperations},
		{HttpMethod: "GET", ResourcePath: "/role/:id/operations", Handler: m.GetRoleOperations},
		{HttpMethod: "GET", ResourcePath: "/user/:id/operations", Handler: m.GetUserOperations},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/operations", Handler: m.SetRoleOperations},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/operations", Handler: m.SetUserOperations},

		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetSystemOperations},
		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/authorities", Handler: m.SetRoleAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/authorities", Handler: m.SetUserAuthorities},
	})
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) GetSystemOperations(context *gin.Context) {
	count, list := c.AuthorityService.GetSystemOperations()
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) GetRoleOperations(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleOperations(roleId)
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetRoleOperations(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	var request []*OperationModel.Operation
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetRoleOperations(roleId, request)

	count, list := c.AuthorityService.GetRoleOperations(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AuthorityMController) GetUserOperations(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserOperations(id); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetUserOperations(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	var request []*OperationModel.Operation
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetUserOperations(userId, request)

	list := c.AuthorityService.GetUserOperations(userId)
	response.Success(context, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) GetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	authType := request.Param(context, "auth_type").DefaultString("menu")

	list, _ := c.AuthorityService.GetRoleAuthorities(roleId, authType)
	response.Page(context, int64(len(list)), list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType := request.Param(context, "auth_type").DefaultString("menu")

	var request []string
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetRoleAuthorities(roleId, authType, request)

	list, _ := c.AuthorityService.GetRoleAuthorities(roleId, authType)
	response.Page(context, int64(len(list)), list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetUserAuthorities(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType := request.Param(context, "auth_type").DefaultString("menu")

	var request []string
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetUserAuthorities(userId, authType, request)

	list := c.AuthorityService.GetUserAuthorities(userId, authType)
	response.Success(context, list)
}

// 根据用户ID获取所有权限的来源
func (c *AuthorityMController) GetUserAuthorities(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType := request.Param(context, "auth_type").DefaultString("menu")

	if data := c.AuthorityService.GetUserAuthorities(id, authType); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}
