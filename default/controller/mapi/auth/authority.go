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

type AuthorityMController struct {
	controller.ResourceController
	AuthorityService *auth.AuthorityService `inject:"authorityService"`
}

var authorityMController *AuthorityMController = &AuthorityMController{}

func init() {
	inject.InjectValue("authorityMController", authorityMController)
}

func (m *AuthorityMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user/menus", Handler: m.GetUserMenus},
		{HttpMethod: "GET", ResourcePath: "/user/menu/:id/buttons", Handler: m.GetUserMenuButtons},
		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetSystemAuthorities},
		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
	})
}

func (c *AuthorityMController) GetUserMenus(context *gin.Context) {
	currentUserId := SecurityUtils.GetCurrentUserId(context)
	menus := c.AuthorityService.GetUserMenuTree(currentUserId)
	if len(menus) > 0 {
		response.Success(context, menus)
	} else {
		response.Success(context, []any{})
	}
}

func (c *AuthorityMController) GetUserMenuButtons(context *gin.Context) {
	menuId, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
	}

	result := c.AuthorityService.GetButtonListByMenuId(SecurityUtils.GetCurrentUserId(context), menuId)
	response.Success(context, result)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) GetSystemAuthorities(context *gin.Context) {
	count, list := c.AuthorityService.GetSystemAuthorities()
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) GetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleAuthorities(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AuthorityMController) GetUserAuthorities(context *gin.Context) {
	id, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserAuthorities(id); data != nil {
		response.Success(context, data)
	} else {
		response.NotFound(context, "")
	}

}
