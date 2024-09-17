package openapi

import (
	"strconv"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/common/service"

	"github.com/gin-gonic/gin"
)

type EventOpenController struct {
	controller.ResourceController
	EventService *service.EventService `inject:"eventService"`
}

var eventOpenController *EventOpenController = &EventOpenController{}

type AdminEventOpenController struct {
	controller.ResourceController
	EventService *service.EventService `inject:"eventService"`
}

var adminEventOpenController *AdminEventOpenController = &AdminEventOpenController{}

func init() {
	inject.InjectValue("eventOpenController", eventOpenController)
	inject.InjectValue("adminEventOpenController", adminEventOpenController)
}

// 组织
func (m *EventOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/events", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/event/:id", Handler: m.GetEvent},
	})
}

// 列表
func (a *EventOpenController) GetList(context *gin.Context) {
	source := request.Param(context, "source").DefaultString("")
	sourceId := request.Param(context, "sourceId").DefaultString("")
	eventType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]interface{})
	if source != "" {
		conds["source"] = source
	}

	if sourceId != "" {
		conds["sourceId"] = sourceId
	}

	if eventType != "" {
		conds["type"] = eventType
	}

	if search != "" {
		conds["search"] = search
	}

	conds["target"] = SecurityUtil.GetCurrentUserId(context)
	conds["scope"] = "PRIVATE"

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.EventService.Find(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}

	eventbus.DispatchEvent("ON_ACCESS_EVENT_CENTER", conds["userId"])
}

// 1.根据id查询节点
func (a *EventOpenController) GetEvent(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	result, err := a.EventService.GetById(id)
	if err != nil {
		response.SystemFail(context, err)
		return
	}

	if result == nil {
		response.NotFound(context, "Not Found")
		return
	}

	tenantId := SecurityUtil.GetCurrentTenantId(context)

	if result.TenantId != tenantId {
		response.NotAllowed(context, "Not Allowed")
		return
	}

	response.Success(context, result)
}

// 组织
func (m *AdminEventOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/events", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/events/head", Handler: m.CheckEvents},
		{HttpMethod: "GET", ResourcePath: "/event/:id", Handler: m.GetEvent},
	})
}

// 列表: Available
func (a *AdminEventOpenController) GetList(context *gin.Context) {
	source := request.Param(context, "source").DefaultString("")
	sourceId := request.Param(context, "sourceId").DefaultString("")
	eventType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]interface{})
	if source != "" {
		conds["source"] = source
	}

	if sourceId != "" {
		conds["sourceId"] = sourceId
	}

	if eventType != "" {
		conds["type"] = eventType
	}

	if search != "" {
		conds["search"] = search
	}

	conds["target"] = SecurityUtil.GetCurrentTenantId(context)
	conds["scope"] = "TENANT"

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.EventService.Find(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}

	eventbus.DispatchEvent("ON_ACCESS_EVENT_CENTER", conds["userId"])
}

// 1.根据id查询节点
func (a *AdminEventOpenController) GetEvent(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	result, err := a.EventService.GetById(id)
	if err != nil {
		response.SystemFail(context, err)
		return
	}

	if result == nil {
		response.NotFound(context, "Not Found")
		return
	}

	tenantId := SecurityUtil.GetCurrentTenantId(context)

	if result.TenantId != tenantId {
		response.NotAllowed(context, "Not Allowed")
		return
	}

	response.Success(context, result)
}

// 列表
func (a *AdminEventOpenController) CheckEvents(context *gin.Context) {
	source := request.Param(context, "source").DefaultString("")
	sourceId := request.Param(context, "sourceId").DefaultString("")
	eventType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")

	var conds = make(map[string]interface{})
	if source != "" {
		conds["source"] = source
	}

	if sourceId != "" {
		conds["sourceId"] = sourceId
	}

	if eventType != "" {
		conds["type"] = eventType
	}

	if search != "" {
		conds["search"] = search
	}

	conds["target"] = SecurityUtil.GetCurrentTenantId(context)
	conds["scope"] = "TENANT"

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if event, err := a.EventService.Check(conds); err == nil {
		response.Success(context, event)
	} else {
		response.Success(context, nil)
	}
}
