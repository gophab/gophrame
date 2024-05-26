package auth

import (
	"strconv"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthService "github.com/gophab/gophrame/default/service/auth"

	"github.com/gin-gonic/gin"
)

type ButtonMController struct {
	controller.ResourceController
	ButtonService     *AuthService.ButtonService     `inject:"buttonService"`
	MenuButtonService *AuthService.MenuButtonService `inject:"menuButtonService"`
}

var buttonMController *ButtonMController = &ButtonMController{}

func init() {
	inject.InjectValue("buttonMController", buttonMController)
}

func (m *ButtonMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/buttons", Handler: m.GetButtons},
		{HttpMethod: "GET", ResourcePath: "/button/:id", Handler: m.GetButton},
		{HttpMethod: "POST", ResourcePath: "/button", Handler: m.CreateButton},
		{HttpMethod: "PUT", ResourcePath: "/button", Handler: m.UpdateButton},
		{HttpMethod: "DELETE", ResourcePath: "/button/:id", Handler: m.DeleteButton},
	})
}

// 1按钮列表
func (s *ButtonMController) GetButtons(context *gin.Context) {
	buttonName := request.Param(context, "buttonName").DefaultString("")
	pageable := query.GetPageable(context)

	count, list := s.ButtonService.List(buttonName, pageable)
	if count > 0 && list != nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, list)
	} else {
		response.Success(context, []any{})
	}
}

func (s *ButtonMController) GetButton(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, err := s.ButtonService.GetById(id); err == nil {
		if res != nil {
			response.Success(c, res)
		} else {
			response.NotFound(c, strconv.FormatInt(id, 10))
		}
	} else {
		response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
	}
}

// 2.按钮新增(store)
func (s *ButtonMController) CreateButton(c *gin.Context) {
	var tmp AuthModel.Button
	if err := c.ShouldBind(&tmp); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if _, err := s.ButtonService.CreateButton(&tmp); err == nil {
		response.Success(c, tmp)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// 5.按钮更新(update)
func (s *ButtonMController) UpdateButton(c *gin.Context) {
	var tmp AuthModel.Button
	if err := c.ShouldBind(&tmp); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	//注意：这里没有实现权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if _, err := s.ButtonService.UpdateButton(&tmp); err == nil {
		response.Success(c, tmp)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}

}

// 6.删除记录
func (u *ButtonMController) DeleteButton(c *gin.Context) {
	buttonId, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	//判断按钮是否被系统菜单引用,如果有,则禁止删除
	if !u.MenuButtonService.GetByButtonId(buttonId) {
		response.FailMessage(c, 400, "该按钮已被菜单引用,无法删除")
	} else {
		if err := u.ButtonService.DeleteButton(buttonId); err == nil {
			response.Success(c, nil)
		} else {
			response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
		}
	}

}
