package mapi

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

var eventMController *EventMController = &EventMController{}

func init() {
	inject.InjectValue("eventMController", eventMController)
}

type EventMController struct {
	controller.ResourceController
	EventService *service.EventService `inject:"eventService"`
}

// 组织
func (m *EventMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/events", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/events/head", Handler: m.CheckEvents},
		{HttpMethod: "GET", ResourcePath: "/event/:id", Handler: m.GetEvent},
	})
}

// 1.根据id查询节点
func (a *EventMController) GetEvent(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if result, _ := a.EventService.GetById(id); result != nil {
		response.Success(context, result)
	} else {
		response.NotFound(context, "")
	}
}

// 列表
func (a *EventMController) GetList(context *gin.Context) {
	source := request.Param(context, "source").DefaultString("")
	sourceId := request.Param(context, "sourceId").DefaultString("")
	eventType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]any)
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

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = "SYSTEM"

	if count, lists, err := a.EventService.Find(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}

	eventbus.DispatchEvent("ON_ACCESS_EVENT_CENTER", conds["userId"])
}

// 列表
func (a *EventMController) CheckEvents(context *gin.Context) {
	source := request.Param(context, "source").DefaultString("")
	sourceId := request.Param(context, "sourceId").DefaultString("")
	eventType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")

	var conds = make(map[string]any)
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

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = "SYSTEM"

	if event, err := a.EventService.Check(conds); err == nil {
		response.Success(context, event)
	} else {
		response.Success(context, nil)
	}
}
