package mapi

import "C"
import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/query"
	SecurityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/wjshen/gophrame/errors"
	"github.com/wjshen/gophrame/service"
	"github.com/wjshen/gophrame/service/auth"
	"github.com/wjshen/gophrame/service/dto"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type UserMController struct {
	controller.ResourceController
	UserService      *service.UserService   `inject:"userService"`
	AuthorityService *auth.AuthorityService `inject:"authorityService"`
}

var userMController *UserMController = &UserMController{}

func init() {
	inject.InjectValue("userMController", userMController)
}

// 用户
func (m *UserMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user", Handler: m.GetCurrentUser},
		{HttpMethod: "GET", ResourcePath: "/users", Handler: m.GetUsers},
		{HttpMethod: "GET", ResourcePath: "/user/:id", Handler: m.GetUser},
		{HttpMethod: "POST", ResourcePath: "/user", Handler: m.AddUser},
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
func (u *UserMController) GetCurrentUser(c *gin.Context) {
	user := SecurityUtil.GetCurrentUser(c)

	if user == nil {
		response.SystemErrorCode(c, errors.ERROR_GET_S_FAIL)
		return
	}

	// menus := make([]string, 1)

	// // 管理员自动增加管理权限
	// if user.Username == "admin" {
	// 	menus = append(menus, "admin")
	// }

	// for _, v := range user.Role {
	// 	exists, err := service.GetRoleService().ExistByID(v.Id)
	// 	if err != nil {
	// 		response.FailCode(c, errors.ERROR_NOT_EXIST)
	// 		return
	// 	}
	// 	if !exists {
	// 		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_NOT_EXIST))
	// 		return
	// 	}

	// 	r, err := service.GetRoleService().Get(v.Id)
	// 	if err != nil {
	// 		response.FailCode(c, errors.ERROR_NOT_EXIST)
	// 		return
	// 	}

	// 	for _, v2 := range r.Menu {
	// 		menus = append(menus, v2.Name)
	// 	}
	// 	menus = util.RemoveRepByMap(menus)
	// }
	//menus := u.AuthorityService.GetUserMenuTree(user.Id)

	response.OK(c, dto.User{User: *user})
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func (u *UserMController) GetUsers(c *gin.Context) {
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
func (u *UserMController) GetUser(c *gin.Context) {
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
func (u *UserMController) AddUser(c *gin.Context) {
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
		response.FailCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	res, err := service.GetUserService().CreateUser(&user)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	err = service.GetUserService().LoadPolicy(res.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EDIT_FAIL)
		return
	}

	response.OK(c, res)
}

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func (u *UserMController) UpdateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
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
		response.SystemErrorCode(c, errors.ERROR_ADD_FAIL)
		return
	}

	exists, err := service.GetUserService().ExistByID(user.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	if result, err := service.GetUserService().UpdateUser(&user); err != nil {
		response.SystemErrorCode(c, errors.ERROR_EDIT_FAIL)
		return
	} else {
		response.OK(c, result)
	}
}

// @Summary   删除用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param  id  path  int true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/users/:id  [DELETE]
func (u *UserMController) DeleteUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).String()

	valid := validation.Validation{}
	valid.MinSize(id, 1, "id").Message("ID不为空")
	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := u.UserService.ExistByID(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	err = u.UserService.Delete(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.OK(c, nil)
}

// @Summary 创建新用户
// @Router /mapi/user [POST]
func (u *UserMController) CreateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		response.FailMessage(c, 400, err.Error())
		return
	}

	if res, err := u.UserService.CreateUser(&user); err != nil {
		response.FailMessage(c, 400, err.Error())
		return
	} else {
		response.Success(c, res)
	}
}
