package openapi

import (
	"encoding/json"

	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	securityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/service"

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
		{HttpMethod: "GET", ResourcePath: "/user/options", Handler: c.GetUserOptions},
		{HttpMethod: "POST", ResourcePath: "/user/options", Handler: c.AddUserOptions},
		{HttpMethod: "PUT", ResourcePath: "/user/options", Handler: c.SetUserOptions},
		{HttpMethod: "DELETE", ResourcePath: "/user/option/:key", Handler: c.RemoveUserOption},
		{HttpMethod: "DELETE", ResourcePath: "/user/options", Handler: c.RemoveUserOptions},
	})
}

func (c *UserOptionOpenController) GetUserOptions(ctx *gin.Context) {
	currentUserId := securityUtil.GetCurrentUserId(ctx)
	if userOptions, err := c.UserOptionService.GetUserOptions(currentUserId); err == nil {
		result := make(map[string]string)
		for name, option := range userOptions.Options {
			result[name] = option.Value
		}
		response.Success(ctx, result)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *UserOptionOpenController) AddUserOptions(ctx *gin.Context) {
	currentUserId := securityUtil.GetCurrentUserId(ctx)
	if body, err := ctx.GetRawData(); err == nil {
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
				response.FailMessage(ctx, 400, err.Error())
				return
			}
		}
		c.GetUserOptions(ctx)
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *UserOptionOpenController) SetUserOptions(ctx *gin.Context) {
	currentUserId := securityUtil.GetCurrentUserId(ctx)
	if body, err := ctx.GetRawData(); err == nil {
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

		if _, err := c.UserOptionService.SetUserOptions(&userOptions); err == nil {
			c.GetUserOptions(ctx)
		} else {
			response.FailMessage(ctx, 400, err.Error())
		}
	} else {
		response.FailMessage(ctx, 400, err.Error())
	}
}

func (c *UserOptionOpenController) RemoveUserOptions(ctx *gin.Context) {
	currentUserId := securityUtil.GetCurrentUserId(ctx)
	if err := c.UserOptionService.RemoveAllUserOptions(currentUserId); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.OK(ctx, nil)
	}
}

func (c *UserOptionOpenController) RemoveUserOption(ctx *gin.Context) {
	currentUserId := securityUtil.GetCurrentUserId(ctx)
	key := ctx.Param("key")
	if res, err := c.UserOptionService.RemoveUserOption(currentUserId, key); err != nil {
		response.FailMessage(ctx, 400, err.Error())
	} else {
		response.OK(ctx, res)
	}
}
