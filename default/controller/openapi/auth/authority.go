package auth

import (
	"strings"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	SecurityUtils "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	"github.com/gophab/gophrame/default/service"
	AuthService "github.com/gophab/gophrame/default/service/auth"

	"github.com/gin-gonic/gin"
)

type AdminAuthorityOpenController struct {
	controller.ResourceController
	AuthorityService *AuthService.AuthorityService `inject:"authorityService"`
	UserServie       *service.UserService          `inject:"userService"`
}

var adminAuthorityOpenController *AdminAuthorityOpenController = &AdminAuthorityOpenController{}

type AuthorityOpenController struct {
	controller.ResourceController
	AuthorityService *AuthService.AuthorityService `inject:"authorityService"`
}

var authorityOpenController *AuthorityOpenController = &AuthorityOpenController{}

func init() {
	inject.InjectValue("authorityOpenController", authorityOpenController)
	inject.InjectValue("adminAuthorityOpenController", adminAuthorityOpenController)
}

func (m *AdminAuthorityOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetConsoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/authorities", Handler: m.SetRoleAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/authorities", Handler: m.SetUserAuthorities},
	})
}

func filter(list []*AuthModel.Operation) []*AuthModel.Operation {
	var result = make([]*AuthModel.Operation, 0)
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
func (c *AdminAuthorityOpenController) GetConsoleAuthorities(context *gin.Context) {
	_, list := c.AuthorityService.GetSystemAuthorities()
	response.Success(context, filter(list))
}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) GetRoleAuthorities(context *gin.Context) {
	roleId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}
	count, list := c.AuthorityService.GetRoleAuthorities(roleId)
	response.Page(context, count, list)
}

// 根据用户ID获取所有权限的来源
func (c *AdminAuthorityOpenController) GetUserAuthorities(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//根据用户ID,查询隶属哪些组织机构
	if data := c.AuthorityService.GetUserAuthorities(userId); data != nil {
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

	var request []*AuthModel.Operation
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetRoleAuthorities(roleId, request)

	count, list := c.AuthorityService.GetRoleAuthorities(roleId)
	response.Page(context, count, list)
}

// 待分配的系统菜单以及挂接的按钮
func (c *AdminAuthorityOpenController) SetUserAuthorities(context *gin.Context) {
	userId, err := request.Param(context, "id").MustString()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	var request []*AuthModel.Operation
	if err := context.ShouldBind(&request); err != nil {
		logger.Warn("数据绑定出错", err.Error())
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	c.AuthorityService.SetUserAuthorities(userId, request)

	list := c.AuthorityService.GetUserAuthorities(userId)
	response.Success(context, list)
}

func (m *AuthorityOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user/menus", Handler: m.GetUserMenus},
		{HttpMethod: "GET", ResourcePath: "/user/menu/:id/buttons", Handler: m.GetUserMenuButtons},
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
