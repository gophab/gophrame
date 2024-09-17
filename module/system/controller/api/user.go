package api

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	auth "github.com/gophab/gophrame/module/authority/service"
	"github.com/gophab/gophrame/module/system/service"
	"github.com/gophab/gophrame/module/system/service/dto"
	"github.com/gophab/gophrame/module/system/service/mapper"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type UserController struct {
	controller.ResourceController
	UserService       *service.UserService       `inject:"userService"`
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
	AuthorityService  *auth.AuthorityService     `inject:"authorityService"`
	UserMapper        *mapper.UserMapper         `inject:"userMapper"`
	SocialUserMapper  *mapper.SocialUserMapper   `inject:"socialUserMapper"`
}

var userController *UserController = &UserController{}

func init() {
	inject.InjectValue("userController", userController)
}

// 用户
func (m *UserController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user", Handler: m.GetCurrentUser},
		{HttpMethod: "GET", ResourcePath: "/users", Handler: m.GetUsers},
		{HttpMethod: "GET", ResourcePath: "/user/:id", Handler: m.GetUser},
		{HttpMethod: "POST", ResourcePath: "/user", Handler: m.CreateUser},
		{HttpMethod: "PUT", ResourcePath: "/user", Handler: m.UpdateUser},
		{HttpMethod: "DELETE", ResourcePath: "/user/:id", Handler: m.DeleteUser},
	})
}

// @Summary   获取登录用户信息
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {"lists":""}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/userInfo  [GET]
func (u *UserController) GetCurrentUser(c *gin.Context) {
	userDetails := SecurityUtil.GetCurrentUser(c)

	if userDetails == nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}
	if userDetails.UserId != nil {
		if user, err := u.UserService.GetById(*userDetails.UserId); err == nil {
			response.Success(c, u.UserMapper.AsDto(user))
		} else {
			response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
		}
	} else {
		if user, err := u.SocialUserService.GetById(*userDetails.SocialId); err == nil {
			response.Success(c, u.SocialUserMapper.AsDto(user))
		} else {
			response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
		}
	}
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func (u *UserController) GetUsers(c *gin.Context) {
	name := request.Param(c, "name").DefaultString("")
	organization := request.Param(c, "organization").DefaultBool(false)

	if organization {
		count, list := u.UserService.GetAllWithOrganization(name, query.GetPageable(c))
		for _, v := range list {
			v.Password = ""
		}
		response.Page(c, count, list)
	} else {
		example := dto.User{}
		example.Login = &name

		count, list := u.UserService.GetAll(&example, query.GetPageable(c))
		for _, v := range list {
			v.Password = ""
		}
		response.Page(c, count, list)
	}
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func (u *UserController) GetUser(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := u.UserService.GetById(id)
	if err != nil {
		response.FailCode(c, errors.ERROR_NOT_EXIST)
		return
	}

	if result == nil {
		response.NotFound(c, "")
		return
	}

	response.Success(c, result)
}

// @Summary   增加用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [POST]
func (u *UserController) CreateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.MaxSize(user.Login, 100, "login").Message("最长为100字符")
	valid.MaxSize(user.Mobile, 20, "mobile").Message("最长为20字符")
	valid.MaxSize(user.Email, 100, "mobile").Message("最长为100字符")
	valid.MaxSize(user.Password, 100, "password").Message("最长为100字符")

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if res, err := service.GetUserService().Create(&user); err == nil {
		response.Success(c, res)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func (u *UserController) UpdateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if user.Id == nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.MinSize(user.Id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(user.Login, 100, "login").Message("最长为100字符")
	valid.MaxSize(user.Mobile, 20, "mobile").Message("最长为20字符")
	valid.MaxSize(user.Email, 100, "email").Message("最长为100字符")
	valid.MaxSize(user.Password, 100, "password").Message("最长为100字符")

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if exists, err := service.GetUserService().ExistByID(*user.Id); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_EXIST_FAIL, err.Error())
		return
	} else if !exists {
		response.NotFound(c, *user.Id)
		return
	}

	if result, err := service.GetUserService().Update(&user); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	} else {
		response.Success(c, result)
	}
}

// @Summary   删除用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param  id  path  int true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/users/:id  [DELETE]
func (u *UserController) DeleteUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).String()

	valid := validation.Validation{}
	valid.MinSize(id, 1, "id").Message("ID不为空")
	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	if exists, err := u.UserService.ExistByID(id); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_EXIST_FAIL, err.Error())
		return
	} else if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	if err := u.UserService.DeleteById(id); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(c, nil)
	}
}
