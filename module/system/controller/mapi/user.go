package mapi

import (
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/core/util/collection"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"
	"github.com/gophab/gophrame/module/system/service/dto"
	"github.com/gophab/gophrame/module/system/service/mapper"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type UserMController struct {
	controller.ResourceController
	UserService   *service.UserService   `inject:"userService"`
	TenantService *service.TenantService `inject:"tenantService"`
	UserMapper    *mapper.UserMapper     `inject:"userMapper"`
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
		{HttpMethod: "PATCH", ResourcePath: "/user/:id", Handler: m.PatchUser},
		{HttpMethod: "DELETE", ResourcePath: "/user/:id", Handler: m.DeleteUser},
		{HttpMethod: "PUT", ResourcePath: "/user/:id/password/reset", Handler: m.ResetUserPassword},
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
	userDetails := SecurityUtil.GetCurrentUser(c)

	if userDetails == nil {
		response.Unauthorized(c, "")
		return
	}
	if userDetails.UserId != nil {
		if user, err := u.UserService.GetById(*userDetails.UserId); err == nil {
			response.Success(c, u.UserMapper.AsDto(user))
		} else {
			response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
		}
	} else {
		response.NotFound(c, "")
	}
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func (u *UserMController) GetUsers(c *gin.Context) {
	search := request.Param(c, "search").DefaultString("")
	id := request.Param(c, "id").DefaultString("")
	name := request.Param(c, "name").DefaultString("")
	login := request.Param(c, "login").DefaultString("")
	mobile := request.Param(c, "mobile").DefaultString("")
	email := request.Param(c, "email").DefaultString("")
	tenantId := request.Param(c, "tenantId").DefaultString("")
	organization := request.Param(c, "organization").DefaultBool(false)

	if organization {
		count, list := u.UserService.GetAllWithOrganization(name, query.GetPageable(c))
		for _, v := range list {
			v.Password = ""
		}
		response.Page(c, count, list)
	} else {
		conds := make(map[string]interface{})
		if search != "" {
			conds["search"] = search
		}
		if id != "" {
			conds["id"] = id
		}
		if name != "" {
			conds["name"] = name
		}
		if login != "" {
			conds["login"] = login
		}
		if mobile != "" {
			conds["mobile"] = mobile
		}
		if email != "" {
			conds["email"] = email
		}
		if tenantId != "" {
			conds["tenantId"] = tenantId
		}

		count, list := u.UserService.Find(conds, query.GetPageable(c))

		tenantIds := collection.MapToSet[string](list, func(i interface{}) string {
			return i.(*domain.User).TenantId
		})

		var tenants = make(map[string]*domain.Tenant)
		if list, err := u.TenantService.GetByIds(tenantIds.AsList()); err == nil {
			for _, item := range list {
				tenants[item.Id] = item
			}
		}
		tenants["SYSTEM"] = &domain.Tenant{
			Id:   "SYSTEM",
			Name: "平台",
		}
		for _, v := range list {
			v.Password = ""
			v.Tenant = tenants[v.TenantId]
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

	if result, err := u.UserService.GetById(id); err == nil {
		if result == nil {
			response.NotFound(c, id)
		} else {
			response.Success(c, result)
		}
	} else {
		response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
	}
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

	// name 不为空
	valid.MaxSize(user.Name, 100, "name").Message("最长为100字符")

	// password 不为空
	valid.MaxSize(*user.PlainPassword, 100, "password").Message("最长为100字符")
	valid.MinSize(*user.PlainPassword, 6, "password").Message("最短为6字符")
	user.Password = user.PlainPassword

	if user.Login != nil {
		if *user.Login == "" {
			user.Login = nil
		} else {
			valid.MaxSize(*user.Login, 100, "login").Message("最长为100字符")
			valid.MinSize(*user.Login, 6, "login").Message("最短为5字符")
		}
	}
	if user.Mobile != nil {
		if *user.Mobile == "" {
			user.Mobile = nil
		} else {
			valid.Check(*user.Mobile, util.NewInternationalTelephoneValidator("mobile")).Message("无效手机号")
		}
	}
	if user.Email != nil {
		if *user.Email == "" {
			user.Email = nil
		} else {
			valid.Email(*user.Email, "email").Message("无效的Email")
		}
	}

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.FailMessage(c, errors.INVALID_PARAMS, valid.Errors[0].Message)
		return
	}

	if res, err := service.GetUserService().Create(&user); err == nil {
		eventbus.PublishEvent("USER_CREATED", res)
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
func (u *UserMController) UpdateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil || user.Id == nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	valid.MinSize(user.Id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(user.Name, 100, "name").Message("最长为100字符")

	if user.Login != nil {
		if *user.Login == "" {
			user.Login = nil
		} else {
			valid.MaxSize(*user.Login, 100, "login").Message("最长为100字符")
			valid.MinSize(*user.Login, 6, "login").Message("最短为5字符")
		}
	}
	if user.Mobile != nil {
		if *user.Mobile == "" {
			user.Mobile = nil
		} else {
			valid.Check(*user.Mobile, util.NewInternationalTelephoneValidator("mobile")).Message("无效手机号")
		}
	}
	if user.Email != nil {
		if *user.Email == "" {
			user.Email = nil
		} else {
			valid.Email(*user.Email, "email").Message("无效的Email")
		}
	}

	if user.PlainPassword != nil {
		if *user.PlainPassword != "" {
			valid.MaxSize(*user.PlainPassword, 100, "password").Message("最长为100字符")
			valid.MinSize(*user.PlainPassword, 6, "password").Message("最短为6字符")
			user.Password = user.PlainPassword
		}
	}

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := service.GetUserService().GetById(*user.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}

	if exists == nil {
		response.NotFound(c, *user.Id)
		return
	}

	if result, err := service.GetUserService().Update(&user); err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
	} else {
		response.Success(c, result)
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
func (u *UserMController) PatchUser(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var params map[string]interface{}
	if err := c.BindJSON(&params); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}

	name := params["name"]
	if name != nil && name.(string) != "" {
		valid.MaxSize(name.(string), 100, "login").Message("最长为100字符")
		valid.MinSize(name.(string), 2, "name").Message("最短为2字符")
	} else {
		delete(params, "name")
	}

	login := params["login"]
	if login != nil && login.(string) != "" {
		valid.MaxSize(login.(string), 100, "login").Message("最长为100字符")
		valid.MinSize(login.(string), 6, "login").Message("最短为5字符")
	} else {
		delete(params, "login")
	}

	mobile := params["mobile"]
	if mobile != nil && mobile.(string) != "" {
		valid.Check(mobile, util.NewInternationalTelephoneValidator("mobile")).Message("无效手机号")
	} else {
		delete(params, "mobile")
	}

	email := params["email"]
	if email != nil && email.(string) != "" {
		valid.Email(email.(string), "email").Message("无效邮箱地址")
	} else {
		delete(params, "email")
	}

	// password := params["password"]
	// if password != nil && password.(string) != "" {
	// 	valid.MaxSize(password.(string), 100, "password").Message("最长为100字符")
	// 	valid.MinSize(password.(string), 6, "password").Message("最短为6字符")
	// }

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, valid.Errors[0].Error())
		return
	}

	var availableFields = []string{"name", "login", "mobile", "email", "admin", "status", "tenantId"}
	var user = make(map[string]interface{})

	for _, k := range availableFields {
		if v, b := params[k]; b {
			user[k] = v
		}
	}

	var result *domain.User
	if result, err = service.GetUserService().PatchAll(id, user); err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
		return
	}

	response.Success(c, result)
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
		response.NotFound(c, id)
		return
	}

	if err = u.UserService.DeleteById(id); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(c, nil)
	}
}

// @Summary 创建新用户
// @Router /mapi/user [POST]
func (u *UserMController) CreateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		response.FailMessage(c, 400, err.Error())
		return
	}

	if res, err := u.UserService.Create(&user); err != nil {
		response.FailMessage(c, 400, err.Error())
		return
	} else {
		response.Success(c, res)
	}
}

func (u *UserMController) ResetUserPassword(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := u.UserService.ExistByID(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, id)
		return
	}

	if b, err := u.UserService.ResetUserPassword(id); err != nil {
		response.SystemError(c, err)
		return
	} else if b {
		response.Success(c, "OK")
		return
	} else {
		response.NotFound(c, "Not Found")
		return
	}
}
