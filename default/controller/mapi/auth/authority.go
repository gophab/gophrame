package auth

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthService "github.com/gophab/gophrame/default/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthorityMController struct {
	controller.ResourceController
	AuthorityService *AuthService.AuthorityService `inject:"authorityService"`
}

var authorityMController *AuthorityMController = &AuthorityMController{}

func init() {
	inject.InjectValue("authorityMController", authorityMController)
}

func (m *AuthorityMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/authorities", Handler: m.GetSystemAuthorities},
		{HttpMethod: "GET", ResourcePath: "/role/:id/authorities", Handler: m.GetRoleAuthorities},
		{HttpMethod: "GET", ResourcePath: "/user/:id/authorities", Handler: m.GetUserAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/role/:id/authorities", Handler: m.SetRoleAuthorities},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/authorities", Handler: m.SetUserAuthorities},
	})
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

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetRoleAuthorities(context *gin.Context) {
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

// 待分配的系统菜单以及挂接的按钮
func (c *AuthorityMController) SetUserAuthorities(context *gin.Context) {
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
