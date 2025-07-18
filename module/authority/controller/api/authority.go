package api

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtils "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/authority/service"

	"github.com/gin-gonic/gin"
)

type AuthorityController struct {
	controller.ResourceController
	AuthorityService *service.AuthorityService `inject:"authorityService"`
}

var authorityController *AuthorityController = &AuthorityController{}

func init() {
	inject.InjectValue("authorityController", authorityController)
}

func (m *AuthorityController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/operations", Handler: m.GetSystemOperations},
		{HttpMethod: "GET", ResourcePath: "/user/menus", Handler: m.GetUserMenus},
		// {HttpMethod: "GET", ResourcePath: "/user/menu/:id/buttons", Handler: m.GetUserMenuButtons},
		{HttpMethod: "GET", ResourcePath: "/role/:id/operations", Handler: m.GetRoleOperations},
		{HttpMethod: "GET", ResourcePath: "/user/:id/operations", Handler: m.GetUserOperations},

		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
	})
}

func (c *AuthorityController) GetUserMenus(context *gin.Context) {
	currentUserId := SecurityUtils.GetCurrentUserId(context)
	menus := c.AuthorityService.GetUserMenuTree(currentUserId)
	if len(menus) > 0 {
		response.Success(context, menus)
	} else {
		response.Success(context, []any{})
	}
}

// func (c *AuthorityController) GetUserMenuButtons(context *gin.Context) {
// 	menuId, err := request.Param(context, "id").MustInt64()
// 	if err != nil {
// 		response.FailCode(context, errors.INVALID_PARAMS)
// 	}

// 	result := c.AuthorityService.GetButtonListByMenuId(SecurityUtils.GetCurrentUserId(context), menuId)
// 	response.Success(context, result)
// }

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityController) GetSystemOperations(context *gin.Context) {
	count, list := c.AuthorityService.GetSystemOperations()
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityController) GetRoleOperations(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleOperations(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AuthorityController) GetUserOperations(context *gin.Context) {
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
func (c *AuthorityController) GetRoleAuthorities(context *gin.Context) {
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
func (c *AuthorityController) GetUserAuthorities(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
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
	if data := c.AuthorityService.GetUserAuthorities(id, authType); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}
