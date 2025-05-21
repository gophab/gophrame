package openapi

import (
	"strings"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	SecurityUtils "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/authority/service"
	"github.com/gophab/gophrame/module/operation/domain"

	"github.com/gin-gonic/gin"
)

type AdminAuthorityOpenController struct {
	controller.ResourceController
	AuthorityService *service.AuthorityService `inject:"authorityService"`
}

var adminAuthorityOpenController *AdminAuthorityOpenController = &AdminAuthorityOpenController{}

type AuthorityOpenController struct {
	controller.ResourceController
	AuthorityService *service.AuthorityService `inject:"authorityService"`
}

var authorityOpenController *AuthorityOpenController = &AuthorityOpenController{}

func init() {
	inject.InjectValue("authorityOpenController", authorityOpenController)
	inject.InjectValue("adminAuthorityOpenController", adminAuthorityOpenController)
}

func (m *AdminAuthorityOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/operations", Handler: m.GetConsoleOperations},
		{HttpMethod: "GET", ResourcePath: "/role/:id/operations", Handler: m.GetRoleOperations},
		{HttpMethod: "GET", ResourcePath: "/user/:id/operations", Handler: m.GetUserOperations},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/operations", Handler: m.SetRoleOperations},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/operations", Handler: m.SetUserOperations},

		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/authorities", Handler: m.SetRoleAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/authorities", Handler: m.SetUserAuthorities},
	})
}

func filter(list []*domain.Operation) []*domain.Operation {
	var result = make([]*domain.Operation, 0)
	if len(list) <= 0 {
		return result
	}

	for _, op := range list {
		if len(op.Children) > 0 {
			// 非叶子节点，看是否所有
			children := filter(op.Children)
			if len(children) > 0 {
				op.Children = children
				result = append(result, op)
			}
		} else {
			if op.Tags == "" {
				result = append(result, op)
			} else if strings.Contains(op.Tags, "console") {
				result = append(result, op)
			}
		}
	}
	return result
}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) GetConsoleOperations(context *gin.Context) {
	_, list := c.AuthorityService.GetSystemOperations()
	response.Success(context, filter(list))
}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) GetRoleOperations(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleOperations(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AdminAuthorityOpenController) GetUserOperations(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserOperations(userId); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) SetRoleOperations(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	var request []*domain.Operation
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetRoleOperations(roleId, request)

	count, list := c.AuthorityService.GetRoleOperations(roleId)
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) SetUserOperations(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	var request []*domain.Operation
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
func (c *AdminAuthorityOpenController) GetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType, err := request.Param(context, "auth_type").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	list, _ := c.AuthorityService.GetRoleAuthorities(roleId, authType)
	response.Page(context, int64(len(list)), list)
}

// 根据用户ID获取所有权限的来源
func (c *AdminAuthorityOpenController) GetUserAuthorities(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType, err := request.Param(context, "auth_type").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserAuthorities(userId, authType); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) SetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType, err := request.Param(context, "auth_type").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

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
func (c *AdminAuthorityOpenController) SetUserAuthorities(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	authType, err := request.Param(context, "auth_type").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

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

func (m *AuthorityOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/operations", Handler: m.GetUserOperations},
		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetUserAuthorities},

		{HttpMethod: "GET", ResourcePath: "/menus", Handler: m.GetUserMenus},
		// {HttpMethod: "GET", ResourcePath: "/user/menu/:id/buttons", Handler: m.GetUserMenuButtons},
	})
}

func (c *AuthorityOpenController) GetUserOperations(context *gin.Context) {
	currentUserId := SecurityUtils.GetCurrentUserId(context)
	operations := c.AuthorityService.GetUserAvailableOperations(currentUserId)
	if len(operations) > 0 {
		response.OK(context, operations)
	} else {
		response.OK(context, []any{})
	}
}

func (c *AuthorityOpenController) GetUserAuthorities(context *gin.Context) {
	authType, err := request.Param(context, "auth_type").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	currentUserId := SecurityUtils.GetCurrentUserId(context)
	authorities := c.AuthorityService.GetUserAvailableAuthorities(currentUserId, authType)
	if len(authorities) > 0 {
		response.OK(context, authorities)
	} else {
		response.OK(context, []any{})
	}
}

func (c *AuthorityOpenController) GetUserMenus(context *gin.Context) {
	currentUserId := SecurityUtils.GetCurrentUserId(context)
	menus := c.AuthorityService.GetUserMenuTree(currentUserId)
	if len(menus) > 0 {
		response.OK(context, menus)
	} else {
		response.OK(context, []any{})
	}
}

// func (c *AuthorityOpenController) GetUserMenuButtons(context *gin.Context) {
// 	menuId, err := request.Param(context, "id").MustInt64()
// 	if err != nil {
// 		response.FailCode(context, errors.INVALID_PARAMS)
// 	}

// 	result := c.AuthorityService.GetButtonListByMenuId(SecurityUtils.GetCurrentUserId(context), menuId)
// 	response.Success(context, result)
// }
