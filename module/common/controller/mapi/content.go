package mapi

import (
	"strconv"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type ContentTemplateMController struct {
	controller.ResourceController
	ContentTemplateService *service.ContentTemplateService `inject:"contentTemplateService"`
}

var contentTemplateMController = &ContentTemplateMController{}

func init() {
	inject.InjectValue("contentTemplateMController", contentTemplateMController)
}

func (c *ContentTemplateMController) AfterInitialize() {
	c.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/content-template/:id", Handler: c.GetContentTemplate},
		{HttpMethod: "GET", ResourcePath: "/content-templates", Handler: c.GetContentTemplates},
		{HttpMethod: "POST", ResourcePath: "/content-template", Handler: c.CreateContentTemplate},
		{HttpMethod: "PUT", ResourcePath: "/content-template", Handler: c.UpdateContentTemplate},
		{HttpMethod: "PATCH", ResourcePath: "/content-template/:id", Handler: c.PatchContentTemplate},
		{HttpMethod: "DELETE", ResourcePath: "/content-template/:id", Handler: c.DeleteContentTemplate},
	})
}

func (c *ContentTemplateMController) GetContentTemplates(ctx *gin.Context) {
	search := request.Param(ctx, "search").DefaultString("")
	typeName := request.Param(ctx, "type").DefaultString("")
	scene := request.Param(ctx, "scene").DefaultString("")

	var conds = make(map[string]any)
	if search != "" {
		conds["search"] = search
	}
	if typeName != "" {
		conds["type"] = typeName
	}
	if scene != "" {
		conds["scene"] = scene
	}

	pageable := query.GetPageable(ctx)

	count, templates, err := c.ContentTemplateService.FindAll(conds, pageable)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	if len(templates) > 0 {
		ctx.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(ctx, templates)
	} else {
		ctx.Header("X-Total-Count", "0")
		response.Success(ctx, []domain.ContentTemplate{})
	}
}

// GET /content-template/:id
// 获取个人会话信息
func (c *ContentTemplateMController) GetContentTemplate(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}

	template, err := c.ContentTemplateService.GetById(id)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	if template == nil {
		response.NotFound(ctx, "")
		return
	}

	response.OK(ctx, template)
}

// POST /template
// 创建新个人会话
func (c *ContentTemplateMController) CreateContentTemplate(ctx *gin.Context) {
	var request domain.ContentTemplate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	if request.Id != "" {
		c.UpdateContentTemplate(ctx)
		return
	}

	request.Status = 1
	request.TenantId = "SYSTEM"
	request.CreatedBy = SecurityUtil.GetCurrentUserId(ctx)

	if template, err := c.ContentTemplateService.CreateContentTemplate(&request); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, template)

		// 操作日志
		eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("CREATE").
			WithTarget("ContentTemplate", template.Id).
			WithContent("${operator.name} 创建了内容模板: ${target.Title} @scence:${target.scene}"))
	}
}

// PUT /template
func (c *ContentTemplateMController) UpdateContentTemplate(ctx *gin.Context) {
	var request domain.ContentTemplate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	if request.Id == "" {
		c.CreateContentTemplate(ctx)
		return
	}

	template, err := c.ContentTemplateService.GetById(request.Id)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	if template == nil {
		response.NotFound(ctx, "")
		return
	}

	template, err = c.ContentTemplateService.UpdateContentTemplate(template)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, template)

	// 操作日志
	eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("UPDATE").
		WithTarget("ContentTemplate", template.Id).
		WithContent("${operator.name} 修改了内容模板: ${target.title} @scence:${target.scene}"))
}

// PUT /template
func (c *ContentTemplateMController) PatchContentTemplate(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}

	var request map[string]any
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	template, err := c.ContentTemplateService.GetById(id)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	if template == nil {
		response.NotFound(ctx, "")
		return
	}

	template, err = c.ContentTemplateService.PatchContentTemplate(id, request)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, template)

	// 操作日志
	eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("UPDATE").
		WithTarget("ContentTemplate", template.Id).
		WithContent("${operator.name} 修改了内容模板: ${target.title} @scence:${target.scene}"))
}

// DELETE /template/:id
// 删除个人会话
func (c *ContentTemplateMController) DeleteContentTemplate(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	template, err := c.ContentTemplateService.GetById(id)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}
	if template == nil {
		response.NotFound(ctx, "")
		return
	}

	err = c.ContentTemplateService.DeleteContentTemplate(template.Id)
	if err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, nil)

	// 操作日志
	eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("DELETE").
		WithTarget("ContentTemplate", template.Id).
		WithContent("${operator.name} 删除了内容模板: ${target.title} @scence:${target.scene}"))
}
