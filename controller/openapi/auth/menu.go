package auth

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"

	securityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/service/auth"

	"github.com/gin-gonic/gin"
)

type MenuOpenController struct {
	controller.ResourceController
	MenuService       *auth.MenuService       `inject:"menuService"`
	MenuButtonService *auth.MenuButtonService `inject:"menuButtonService"`
	AuthorityService  *auth.AuthorityService  `inject:"authorityService"`
}

var menuOpenController *MenuOpenController = &MenuOpenController{}

func init() {
	inject.InjectValue("menuOpenController", menuOpenController)
}

func (m *MenuOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/menus", Handler: m.GetCurrentUserMenus},
		{HttpMethod: "GET", ResourcePath: "/menus/:id", Handler: m.GetSubMenus},
		{HttpMethod: "GET", ResourcePath: "/menu/:id", Handler: m.GetMenu},
	})
}

// @Summary   获取登录用户信息
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {"lists":""}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/userInfo  [GET]
func (u *MenuOpenController) GetCurrentUserMenus(c *gin.Context) {
	userId := securityUtil.GetCurrentUserId(c)

	if userId == "" {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	menus := u.AuthorityService.GetUserMenuTree(userId)
	response.OK(c, menus)
}

func (u *MenuOpenController) GetMenu(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	userId := securityUtil.GetCurrentUserId(c)
	if userId == "" {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	menus := u.AuthorityService.GetUserMenus(userId)
	for _, menu := range menus {
		if menu.Id == id {
			response.Success(c, menu)
			return
		}
	}
	response.NotFound(c, "")
}

func (u *MenuOpenController) GetSubMenus(c *gin.Context) {
	fid, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	userId := securityUtil.GetCurrentUserId(c)
	if userId == "" {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	menus := u.AuthorityService.GetUserMenus(userId)
	for _, menu := range menus {
		if menu.Id == fid {
			if res, err := u.MenuService.GetByFid(fid); err == nil {
				response.Success(c, res)
			} else {
				response.FailMessage(c, 400, err.Error())
			}
			return
		}
	}
	response.NotFound(c, "")
}
