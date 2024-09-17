package mapi

import (
	"strings"

	"github.com/gophab/gophrame/module/slink/config"
	"github.com/gophab/gophrame/module/slink/service"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gin-gonic/gin"
)

type ShortLinkMController struct {
	controller.ResourceController
	ShortLinkService *service.ShortLinkService `inject:"shortLinkService"`
}

var shortLinkMController = &ShortLinkMController{}

func init() {
	inject.InjectValue("shortLinkMController", shortLinkMController)
}

func (c *ShortLinkMController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/short-link/:id", Handler: c.GetShortLink},
		{HttpMethod: "GET", ResourcePath: "/short-links", Handler: c.GetShortLinks},
		{HttpMethod: "POST", ResourcePath: "/short-link", Handler: c.CreateShortLink},
		{HttpMethod: "PATCH", ResourcePath: "/short-link/:id", Handler: c.PatchShortLink},
		{HttpMethod: "DELETE", ResourcePath: "/short-link/:id", Handler: c.DeleteShortLink},
	})
}

func (c *ShortLinkMController) GetShortLink(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.ShortLinkService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}

	result, err = c.ShortLinkService.GetByKey(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}
}

func (c *ShortLinkMController) GetShortLinks(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.ShortLinkService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}

	result, err = c.ShortLinkService.GetByKey(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}
}

func (c *ShortLinkMController) CreateShortLink(ctx *gin.Context) {
	var request struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.ShortLinkService.GenerateShortLink(request.Name, request.Url, config.Setting.Expired)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		root, _ := strings.CutSuffix(config.Setting.BaseUrl, "/")
		response.Success(ctx, root+"/"+result.Key)
		return
	}
}

func (c *ShortLinkMController) PatchShortLink(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.ShortLinkService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}

	result, err = c.ShortLinkService.GetByKey(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result != nil {
		response.Success(ctx, result)
		return
	}
}

func (c *ShortLinkMController) DeleteShortLink(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	err = c.ShortLinkService.DeleteById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	err = c.ShortLinkService.DeleteByKey(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	response.Success(ctx, nil)
}
