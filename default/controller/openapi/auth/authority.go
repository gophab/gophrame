package auth

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtils "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/default/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthorityOpenController struct {
	controller.ResourceController
	AuthorityService *auth.AuthorityService `inject:"authorityService"`
}

var authorityOpenController *AuthorityOpenController = &AuthorityOpenController{}

func init() {
	inject.InjectValue("authorityOpenController", authorityOpenController)
}

func (m *AuthorityOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user/menus", Handler: m.GetUserMenus},
		{HttpMethod: "GET", ResourcePath: "/user/menu/:id/buttons", Handler: m.GetUserMenuButtons},
		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetSystemAuthorities},
		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
	})
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

func (c *AuthorityOpenController) GetUserMenuButtons(context *gin.Context) {
	menuId, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
	}

	result := c.AuthorityService.GetButtonListByMenuId(SecurityUtils.GetCurrentUserId(context), menuId)
	response.Success(context, result)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityOpenController) GetSystemAuthorities(context *gin.Context) {
	count, list := c.AuthorityService.GetSystemAuthorities()
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityOpenController) GetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleAuthorities(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AuthorityOpenController) GetUserAuthorities(context *gin.Context) {
	currentUserId := SecurityUtils.GetCurrentUserId(context)

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserAuthorities(currentUserId); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}
