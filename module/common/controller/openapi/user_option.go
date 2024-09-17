package openapi

import (
	"encoding/json"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"

	"github.com/gin-gonic/gin"
)

type UserOptionOpenController struct {
	controller.ResourceController
	UserOptionService *service.UserOptionService `inject:"userOptionService"`
}

var userOptionOpenController = &UserOptionOpenController{}

func init() {
	inject.InjectValue("userOptionOpenController", userOptionOpenController)
}

func (c *UserOptionOpenController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/options", Handler: c.GetUserOptions},
		{HttpMethod: "POST", ResourcePath: "/options", Handler: c.AddUserOptions},
		{HttpMethod: "PUT", ResourcePath: "/options", Handler: c.SetUserOptions},
		{HttpMethod: "DELETE", ResourcePath: "/option/:key", Handler: c.RemoveUserOption},
		{HttpMethod: "DELETE", ResourcePath: "/options", Handler: c.RemoveUserOptions},
	})
}

func (c *UserOptionOpenController) GetUserOptions(ctx *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	if userOptions, err := c.UserOptionService.GetUserOptions(currentUserId); err == nil {
		result := make(map[string]string)
		for name, option := range userOptions.Options {
			result[name] = option.Value
		}
		response.Success(ctx, result)
	} else {
		response.SystemErrorMessage(ctx, errors.ERROR_GET_S_FAIL, err.Error())
	}
}

func (c *UserOptionOpenController) AddUserOptions(ctx *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	body, err := ctx.GetRawData()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	var data map[string]string
	_ = json.Unmarshal(body, &data)
	for k, v := range data {
		if _, err := c.UserOptionService.AddUserOption(&domain.UserOption{
			UserId: currentUserId,
			Option: domain.Option{
				Name:      k,
				Value:     v,
				ValueType: "STRING",
			},
		}); err != nil {
			response.SystemErrorMessage(ctx, errors.ERROR_CREATE_FAIL, err.Error())
			return
		}
	}
	c.GetUserOptions(ctx)
}

func (c *UserOptionOpenController) SetUserOptions(ctx *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	body, err := ctx.GetRawData()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	var userOptions = domain.UserOptions{
		UserId: currentUserId,
	}

	var data map[string]string
	_ = json.Unmarshal(body, &data)
	for k, v := range data {
		userOptions.Options[k] = domain.UserOption{
			UserId: currentUserId,
			Option: domain.Option{
				Name:      k,
				Value:     v,
				ValueType: "STRING",
			},
		}
	}

	if _, err := c.UserOptionService.SetUserOptions(&userOptions); err != nil {
		response.SystemErrorMessage(ctx, errors.ERROR_UPDATE_FAIL, err.Error())
		return
	}

	c.GetUserOptions(ctx)
}

func (c *UserOptionOpenController) RemoveUserOptions(ctx *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	if err := c.UserOptionService.RemoveAllUserOptions(currentUserId); err != nil {
		response.SystemErrorMessage(ctx, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *UserOptionOpenController) RemoveUserOption(ctx *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	key := ctx.Param("key")
	if res, err := c.UserOptionService.RemoveUserOption(currentUserId, key); err != nil {
		response.SystemErrorMessage(ctx, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(ctx, res)
	}
}
