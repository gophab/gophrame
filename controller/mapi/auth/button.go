package auth

import (
	"strconv"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"

	AuthModel "github.com/wjshen/gophrame/domain/auth"
	"github.com/wjshen/gophrame/service/auth"

	"github.com/gin-gonic/gin"
)

type ButtonMController struct {
	controller.ResourceController
	ButtonService     *auth.ButtonService     `inject:"buttonService"`
	MenuButtonService *auth.MenuButtonService `inject:"menuButtonService"`
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

func (s *ButtonMController) GetButton(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if res, err := s.ButtonService.GetById(id); err == nil {
		response.Success(context, res)
	} else {
		response.FailMessage(context, 400, err.Error())
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
		response.Fail(c, 500, err.Error())
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
		response.Fail(c, 500, err.Error())
	}

}

// 6.删除记录
func (u *ButtonMController) DeleteButton(context *gin.Context) {
	buttonId, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	//判断按钮是否被系统菜单引用,如果有,则禁止删除
	if !u.MenuButtonService.GetByButtonId(buttonId) {
		response.Fail(context, 400, "该按钮已被菜单引用,无法删除")
	} else {
		if err := u.ButtonService.DeleteButton(buttonId); err == nil {
			response.Success(context, nil)
		} else {
			response.Fail(context, 400, err.Error())
		}
	}

}
