package mapi

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

type MenuMController struct {
	controller.ResourceController
	MenuService *service.MenuService `inject:"menuService"`
}

var menuMController *MenuMController = &MenuMController{}

func init() {
	inject.InjectValue("menuMController", menuMController)
}

func (m *MenuMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/menus", Handler: m.GetMenus},
		{HttpMethod: "GET", ResourcePath: "/menus/:id", Handler: m.GetSubMenus},
		{HttpMethod: "GET", ResourcePath: "/menu/:id", Handler: m.GetMenu},
		{HttpMethod: "POST", ResourcePath: "/menu", Handler: m.CreateMenu},
		{HttpMethod: "PUT", ResourcePath: "/menu", Handler: m.UpdateMenu},
		{HttpMethod: "DELETE", ResourcePath: "/menu/:id", Handler: m.DeleteMenu},
	})
}

func (m *MenuMController) CreateMenu(c *gin.Context) {
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

func (m *MenuMController) UpdateMenu(c *gin.Context) {
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

func (m *MenuMController) DeleteMenu(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
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

func (m *MenuMController) GetMenu(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, _ := m.MenuService.GetById(id); res == nil {
		response.NotFound(c, "Not Found")
	} else {
		response.Success(c, res)
	}
}

func (m *MenuMController) GetMenus(c *gin.Context) {
	fid := request.Param(c, "fid").DefaultInt64(-1)
	title := request.Param(c, "title").DefaultString("")
	buttons := request.Param(c, "buttons").DefaultBool(false)
	tree := request.Param(c, "tree").DefaultBool(false)
	pageable := query.GetPageable(c)

	if buttons { // menu with buttons
		count, result := m.MenuService.ListWithButtons(fid, title, pageable)
		if count != 0 {
			c.Header("X-Total-Count", strconv.FormatInt(count, 10))
			if tree {
				result = m.MenuService.MakeTree(result)
			}
			response.Success(c, result)
			return
		}
	} else {
		count, result := m.MenuService.List(fid, title, pageable)
		if count != 0 {
			c.Header("X-Total-Count", strconv.FormatInt(count, 10))
			if tree {
				result = m.MenuService.MakeTree(result)
			}
			response.Success(c, result)
			return
		}
	}
	response.Success(c, []any{})
}

func (m *MenuMController) GetSubMenus(c *gin.Context) {
	fid, err := request.Param(c, "id").MustInt64()
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
