package api

import (
	"strconv"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/operation/domain"
	"github.com/gophab/gophrame/module/operation/service"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	controller.ResourceController
	MenuService *service.MenuService `inject:"menuService"`
}

var menuController *MenuController = &MenuController{}

func init() {
	inject.InjectValue("menuController", menuController)
}

func (m *MenuController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/menus", Handler: m.GetMenus},
		{HttpMethod: "GET", ResourcePath: "/menus/:id", Handler: m.GetSubMenus},
		{HttpMethod: "GET", ResourcePath: "/menu/:id", Handler: m.GetMenu},
		{HttpMethod: "POST", ResourcePath: "/menu", Handler: m.CreateMenu},
		{HttpMethod: "PUT", ResourcePath: "/menu", Handler: m.UpdateMenu},
		{HttpMethod: "DELETE", ResourcePath: "/menu/:id", Handler: m.DeleteMenu},
	})
}

func (m *MenuController) CreateMenu(c *gin.Context) {
	var tmp domain.Menu
	if err := c.ShouldBind(&tmp); err == nil {
		if b, err := m.MenuService.CreateMenu(&tmp); b {
			response.Success(c, tmp)
		} else {
			response.FailMessage(c, 400, err.Error())
		}
	} else {
		logger.Warn("MenuRepository 数据绑定出错", err.Error())
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *MenuController) UpdateMenu(c *gin.Context) {
	var tmp domain.Menu
	if err := c.ShouldBind(&tmp); err == nil {
		if b, err := m.MenuService.UpdateMenu(&tmp); b {
			response.Success(c, tmp)
		} else {
			response.FailMessage(c, 400, err.Error())
		}
	} else {
		logger.Warn("MenuRepository 数据绑定出错", err.Error())
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *MenuController) DeleteMenu(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if b, err := m.MenuService.DeleteMenu(id); b {
		response.Success(c, "")
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

func (m *MenuController) GetMenu(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, _ := m.MenuService.GetById(id); res != nil {
		response.NotFound(c, "")
	} else {
		response.Success(c, res)
	}
}

func (m *MenuController) GetMenus(c *gin.Context) {
	fid := request.Param(c, "fid").DefaultString("")
	title := request.Param(c, "title").DefaultString("")
	buttons := request.Param(c, "buttons").DefaultBool(false)
	pageable := query.GetPageable(c)

	if buttons { // menu with buttons
		count, result := m.MenuService.ListWithButtons(fid, title, pageable)
		if count > 0 {
			c.Header("X-Total-Count", strconv.FormatInt(count, 10))
			response.Success(c, result)
			return
		}
	} else {
		count, result := m.MenuService.List(fid, title, pageable)
		if count > 0 {
			c.Header("X-Total-Count", strconv.FormatInt(count, 10))
			response.Success(c, result)
			return
		}
	}
	response.Success(c, []any{})
}

func (m *MenuController) GetSubMenus(c *gin.Context) {
	fid, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, err := m.MenuService.GetByFid(fid); err == nil {
		response.Success(c, res)
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}
