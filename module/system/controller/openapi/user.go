package openapi

import (
	"bufio"
	"strings"

	"github.com/gophab/gophrame/core/controller"
	EmailCode "github.com/gophab/gophrame/core/email/code"
	"github.com/gophab/gophrame/core/excel"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	SmsCode "github.com/gophab/gophrame/core/sms/code"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"
	"github.com/gophab/gophrame/module/system/service/dto"
	"github.com/gophab/gophrame/module/system/service/mapper"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type UserOpenController struct {
	controller.ResourceController
	UserService        *service.UserService          `inject:"userService"`
	SocialUserService  *service.SocialUserService    `inject:"socialUserService"`
	InviteCodeService  *service.InviteCodeService    `inject:"inviteCodeService"`
	UserMapper         *mapper.UserMapper            `inject:"userMapper"`
	SocialUserMapper   *mapper.SocialUserMapper      `inject:"socialUserMapper"`
	SmsCodeValidator   *SmsCode.SmsCodeValidator     `inject:"smsCodeValidator"`
	EmailCodeValidator *EmailCode.EmailCodeValidator `inject:"emailCodeValidator"`
}

var userOpenController *UserOpenController = &UserOpenController{}

type AdminUserOpenController struct {
	controller.ResourceController
	UserService       *service.UserService       `inject:"userService"`
	TenantService     *service.TenantService     `inject:"tenantService"`
	SocialUserService *service.SocialUserService `inject:"socialUserService"`
	InviteCodeService *service.InviteCodeService `inject:"inviteCodeService"`
	UserMapper        *mapper.UserMapper         `inject:"userMapper"`
	SocialUserMapper  *mapper.SocialUserMapper   `inject:"socialUserMapper"`
}

var adminUserOpenController *AdminUserOpenController = &AdminUserOpenController{}

func init() {
	inject.InjectValue("userOpenController", userOpenController)
	inject.InjectValue("adminUserOpenController", adminUserOpenController)

}

// 用户
func (m *UserOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user", Handler: m.GetCurrentUser},
		{HttpMethod: "GET", ResourcePath: "/user/:id", Handler: m.GetUser},
		{HttpMethod: "GET", ResourcePath: "/user/invite-code", Handler: m.GetUserInviteCode},
		{HttpMethod: "PUT", ResourcePath: "/user", Handler: m.UpdateUser},
		{HttpMethod: "PATCH", ResourcePath: "/user", Handler: m.PatchUser},
		{HttpMethod: "PUT", ResourcePath: "/user/mobile", Handler: m.ChangeUserMobile},
		{HttpMethod: "PUT", ResourcePath: "/user/email", Handler: m.ChangeUserEmail},
		{HttpMethod: "DELETE", ResourcePath: "/user", Handler: m.DeleteUser},
		{HttpMethod: "PUT", ResourcePath: "/user/password", Handler: m.ChangeUserPassword},
	})
}

// @Summary   获取登录用户信息
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {"lists":""}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/userInfo  [GET]
func (u *UserOpenController) GetCurrentUser(c *gin.Context) {
	userId := SecurityUtil.GetCurrentUserId(c)

	if userId == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	if strings.HasPrefix(userId, "sns_") {
		// 社交账户登录
		if user, err := u.SocialUserService.GetById(userId); err == nil {
			response.Success(c, u.SocialUserMapper.AsDto(user))
		} else {
			response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
		}
	} else {
		if user, err := u.UserService.GetById(userId); err == nil {
			response.Success(c, u.UserMapper.AsDto(user))
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
func (u *UserOpenController) GetUser(c *gin.Context) {
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

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func (u *UserOpenController) UpdateUser(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil || user.Id == nil {
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
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	exists, err := service.GetUserService().ExistByID(*user.Id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if !exists {
		response.NotFound(c, errors.GetErrorMessage(errors.ERROR_EXIST_FAIL))
		return
	}

	if result, err := service.GetUserService().Update(&user); err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
		return
	} else {
		response.OK(c, result)
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
func (u *UserOpenController) PatchUser(c *gin.Context) {
	var params map[string]interface{}
	if err := c.BindJSON(&params); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := service.GetUserService().PatchAll(SecurityUtil.GetCurrentUserId(c), params)
	if err != nil {
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
// @Router /openapi/user  [DELETE]
func (u *UserOpenController) DeleteUser(c *gin.Context) {
	id := SecurityUtil.GetCurrentUserId(c)

	valid := validation.Validation{}
	valid.MinSize(id, 1, "id").Message("ID不为空")
	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.INVALID_PARAMS)
		return
	}

	err := u.UserService.DeleteById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.OK(c, nil)
}

/**
 * GET /openapi/user/invite-code
 */
func (u *UserOpenController) GetUserInviteCode(c *gin.Context) {
	channel := request.Param(c, "channel").DefaultString("INVITE_REGISTER")
	currentUserId := SecurityUtil.GetCurrentUserId(c)
	result, err := u.InviteCodeService.GetUserInviteCode(currentUserId, channel)
	if err != nil {
		response.FailMessage(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

/**
 * GET /openapi/user/invite-code
 */
type ChangePasswordForm struct {
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}

func (u *UserOpenController) ChangeUserPassword(c *gin.Context) {
	var request ChangePasswordForm
	if err := c.BindJSON(&request); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if request.Password == "" || request.OldPassword == "" {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if b, err := u.UserService.ChangeUserPassword(SecurityUtil.GetCurrentUserId(c), request.OldPassword, request.Password); b {
		response.Success(c, "OK")
	} else if err != nil {
		response.SystemFail(c, err)
	} else {
		response.NotFound(c, "Not Found")
	}
}

type ChangeForm struct {
	Target string `json:"target"`
	Code   string `json:"code"`
}

func (u *UserOpenController) ChangeUserMobile(c *gin.Context) {
	var request ChangeForm
	if err := c.BindJSON(&request); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if request.Target == "" || request.Code == "" {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if !u.SmsCodeValidator.CheckCode(u.SmsCodeValidator, request.Target, "bind", request.Code) {
		response.Forbidden(c, "Forbidden")
		return
	}

	if _, err := u.UserService.Patch(SecurityUtil.GetCurrentUserId(c), "mobile", request.Target); err == nil {
		response.Success(c, nil)
	} else {
		response.SystemFail(c, err)
	}
}

func (u *UserOpenController) ChangeUserEmail(c *gin.Context) {
	var request ChangeForm
	if err := c.BindJSON(&request); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if request.Target == "" || request.Code == "" {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if !u.EmailCodeValidator.CheckCode(u.EmailCodeValidator, request.Target, "bind", request.Code) {
		response.Forbidden(c, "Forbidden")
		return
	}

	if _, err := u.UserService.Patch(SecurityUtil.GetCurrentUserId(c), "email", request.Target); err == nil {
		response.Success(c, nil)
	} else {
		response.SystemFail(c, err)
	}
}

// 用户
func (m *AdminUserOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/user", Handler: m.GetCurrentUser},
		{HttpMethod: "GET", ResourcePath: "/user/:id", Handler: m.GetUser},
		{HttpMethod: "GET", ResourcePath: "/users", Handler: m.GetUsers},
		{HttpMethod: "POST", ResourcePath: "/user", Handler: m.CreateUser},
		{HttpMethod: "POST", ResourcePath: "/users", Handler: m.CreateUsers},
		{HttpMethod: "POST", ResourcePath: "/users/import", Handler: m.ImportUsers},
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
func (u *AdminUserOpenController) GetCurrentUser(c *gin.Context) {
	userId := SecurityUtil.GetCurrentUserId(c)

	if userId == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	if strings.HasPrefix(userId, "sns_") {
		// 社交账户登录
		if user, err := u.SocialUserService.GetById(userId); err == nil {
			response.Success(c, u.SocialUserMapper.AsDto(user))
		} else {
			response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
		}
	} else {
		if user, err := u.UserService.GetById(userId); err == nil {
			response.Success(c, u.UserMapper.AsDto(user))
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
func (u *AdminUserOpenController) GetUsers(c *gin.Context) {
	search := request.Param(c, "search").DefaultString("")
	id := request.Param(c, "id").DefaultString("")
	name := request.Param(c, "name").DefaultString("")
	login := request.Param(c, "login").DefaultString("")
	mobile := request.Param(c, "mobile").DefaultString("")
	email := request.Param(c, "email").DefaultString("")
	organization := request.Param(c, "organization").DefaultBool(false)
	tenantId := SecurityUtil.GetCurrentTenantId(c)

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

	if organization {
		count, list := u.UserService.GetAllWithOrganization(name, query.GetPageable(c))
		u.UserService.LoadUsersRoles(list)
		response.Page(c, count, list)
	} else {
		count, list := u.UserService.Find(conds, query.GetPageable(c))
		u.UserService.LoadUsersRoles(list)
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
func (u *AdminUserOpenController) GetUser(c *gin.Context) {
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

	u.UserService.LoadUserRoles(result)

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
func (u *AdminUserOpenController) CreateUser(c *gin.Context) {
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
			valid.MinSize(*user.Login, 5, "login").Message("最短为5字符")
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
		response.FailMessage(c, errors.ERROR_CREATE_FAIL, valid.Errors[0].Message)
		return
	}

	user.TenantId = util.StringAddr(SecurityUtil.GetCurrentTenantId(c))

	res, err := service.GetUserService().Create(&user)
	if err != nil {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
		return
	}

	u.UserService.LoadUserRoles(res)

	response.OK(c, res)
}

// @Summary   增加用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [POST]
func (u *AdminUserOpenController) CreateUsers(c *gin.Context) {
	var users []dto.User
	if err := c.BindJSON(&users); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}
	result := make([]*domain.User, 0)
	for _, user := range users {
		var skip = false
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
				valid.MinSize(*user.Login, 5, "login").Message("最短为5字符")
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
			skip = true
		}

		if !skip {
			user.TenantId = util.StringAddr(SecurityUtil.GetCurrentTenantId(c))
			res, err := service.GetUserService().Create(&user)
			if err != nil {
				logger.Error("Create user error: ", err.Error())
			} else {
				result = append(result, res)
			}
		}
	}

	u.UserService.LoadUsersRoles(result)

	response.OK(c, result)
}

// @Summary   增加用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [POST]
func (u *AdminUserOpenController) ImportUsers(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.FailMessage(c, 400, "接收文件失败")
		return
	}

	tenantId := SecurityUtil.GetCurrentTenantId(c)

	// 异步执行导入
	go func() {
		defer func() {
			file.Close()
		}()

		excel.NewImporter().ReadReader(bufio.NewReader(file), func(rows []map[string]string) {
			for _, row := range rows {
				var skip = false
				// 处理每条记录
				var user dto.User
				user.Name = util.Nullable(util.StringAddr(row["name"]))
				user.Login = util.Nullable(util.StringAddr(row["login"]))
				user.Email = util.Nullable(util.StringAddr(row["email"]))
				user.Mobile = util.Nullable(util.StringAddr(row["mobile"]))
				user.Password = util.Nullable(util.StringAddr(row["password"]))
				user.PlainPassword = util.Nullable(util.StringAddr("!12345678!"))

				if user.Name == nil && user.Mobile == nil && user.Email == nil {
					skip = true
				}

				valid := validation.Validation{}

				valid.MaxSize(user.Name, 100, "name").Message("最长为100字符")

				if user.Login != nil && *user.Login != "" {
					valid.Min(*user.Login, 5, "login").Message("至少5个字符")
					valid.AlphaDash(*user.Login, "login").Message("账号格式错误")
				} else {
					user.Login = nil
				}

				if user.Mobile != nil && *user.Mobile != "" {
					valid.Check(*user.Mobile, util.NewInternationalTelephoneValidator("mobile")).Message("无效手机号")
				} else {
					user.Mobile = nil
				}

				if user.Email != nil && *user.Email != "" {
					valid.Email(user.Email, "email").Message("无效的邮箱地址")
				} else {
					user.Email = nil
				}

				if user.Password != nil && *user.Password != "" {
					valid.MaxSize(*user.Password, 100, "password").Message("最长为100字符")
					valid.MinSize(*user.Password, 6, "password").Message("最短为6字符")
				} else {
					user.Password = nil
				}

				if valid.HasErrors() {
					skip = true
				}

				if !skip {
					user.TenantId = util.StringAddr(tenantId)
					u.UserService.Create(&user)
				}
			}
		}, 50)
	}()

	response.OK(c, nil)
}

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func (u *AdminUserOpenController) UpdateUser(c *gin.Context) {
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
			valid.MinSize(*user.Login, 5, "login").Message("最短为5字符")
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
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
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

	if exists.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	result, err := service.GetUserService().Update(&user)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
		return
	}

	u.UserService.LoadUserRoles(result)
	response.OK(c, result)
}

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func (u *AdminUserOpenController) PatchUser(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil || id == "" {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := service.GetUserService().GetById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}

	if exists == nil {
		response.NotFound(c, "Not Found")
		return
	}

	if exists.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	var user = make(map[string]interface{})

	var params domain.User
	if err := c.BindJSON(&params); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	valid := validation.Validation{}

	name := params.Name
	if name != nil && *name != "" {
		valid.MaxSize(*name, 100, "name").Message("最长为100字符")
		valid.MinSize(*name, 2, "name").Message("最短为2字符")

		user["name"] = *name
	}

	login := params.Login
	if login != nil && *login != "" {
		valid.MaxSize(*login, 100, "login").Message("最长为100字符")
		valid.MinSize(*login, 6, "login").Message("最短为5字符")

		user["login"] = *login
	}

	mobile := params.Mobile
	if mobile != nil && *mobile != "" {
		valid.Check(*mobile, util.NewInternationalTelephoneValidator("mobile")).Message("无效手机号")

		user["mobile"] = *mobile
	}

	email := params.Email
	if email != nil && *email != "" {
		valid.Email(*email, "email").Message("无效邮箱地址")

		user["email"] = *email
	}

	if params.Avatar != nil {
		user["avatar"] = *params.Avatar
	}

	user["admin"] = params.Admin

	if params.Status != nil {
		user["status"] = *params.Status
	}

	if params.Roles != nil {
		user["roles"] = params.Roles
	}

	// password := params["password"]
	// if password != nil && password.(string) != "" {
	// 	valid.MaxSize(password.(string), 100, "password").Message("最长为100字符")
	// 	valid.MinSize(password.(string), 6, "password").Message("最短为6字符")
	// }

	if valid.HasErrors() {
		logger.MarkErrors(valid.Errors)
		response.SystemErrorCode(c, errors.ERROR_CREATE_FAIL)
		return
	}

	// var availableFields = []string{"name", "avatar", "login", "mobile", "email", "admin", "status"}

	result, err := service.GetUserService().PatchAll(id, user)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_UPDATE_FAIL)
		return
	}

	u.UserService.LoadUserRoles(result)
	response.Success(c, result)
}

// @Summary   删除用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param  id  path  int true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/users/:id  [DELETE]
func (u *AdminUserOpenController) DeleteUser(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exists, err := service.GetUserService().GetById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}

	if exists == nil {
		response.NotFound(c, "Not Found")
		return
	}

	if exists.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	err = u.UserService.DeleteById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_DELETE_FAIL)
		return
	}

	response.OK(c, nil)
}

func (u *AdminUserOpenController) ResetUserPassword(c *gin.Context) {
	id, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	exist, err := u.UserService.GetById(id)
	if err != nil {
		response.SystemErrorCode(c, errors.ERROR_EXIST_FAIL)
		return
	}
	if exist == nil {
		response.NotFound(c, id)
		return
	}

	if exist.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	if exist.Id == SecurityUtil.GetCurrentUserId(c) {
		response.NotAllowed(c, "Not Allowed")
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
