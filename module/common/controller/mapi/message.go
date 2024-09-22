package mapi

import (
	"strconv"
	"time"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"

	"github.com/gin-gonic/gin"
)

var messageMController *MessageMController = &MessageMController{}

func init() {
	inject.InjectValue("messageMController", messageMController)
}

type MessageMController struct {
	controller.ResourceController
	MessageService *service.MessageService `inject:"messageService"`
}

// 组织
func (m *MessageMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/messages", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/messages/managed", Handler: m.GetManagedList},
		{HttpMethod: "GET", ResourcePath: "/message/:id", Handler: m.GetMessage},
		{HttpMethod: "POST", ResourcePath: "/message", Handler: m.CreateMessage},
		{HttpMethod: "PUT", ResourcePath: "/message", Handler: m.UpdateMessage},
		{HttpMethod: "PATCH", ResourcePath: "/message/:id", Handler: m.PatchMessage},
		{HttpMethod: "DELETE", ResourcePath: "/message/:id", Handler: m.DeleteMessage},
	})
}

// 1.根据id查询节点
func (a *MessageMController) GetMessage(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	show := request.Param(context, "show").DefaultBool(false)
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	if result, _ := a.MessageService.GetById(id); result != nil {
		response.Success(context, result)

		if show {
			eventbus.DispatchEvent("SYSTEM_MESSAGE_VIEWED", result, SecurityUtil.GetCurrentUserId(context))
		}

	} else {
		response.NotFound(context, "")
	}
}

// 列表
func (a *MessageMController) GetList(context *gin.Context) {
	messageFrom := request.Param(context, "from").DefaultString("")
	messageType := request.Param(context, "type").DefaultString("")
	messageScope := request.Param(context, "scope").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]interface{})
	if messageFrom != "" {
		conds["from"] = messageFrom
	}

	if messageType != "" {
		conds["type"] = messageType
	}

	if messageScope != "" {
		conds["scope"] = messageScope
	}

	if search != "" {
		conds["search"] = search
	}

	conds["status"] = 1
	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = "SYSTEM"

	if count, lists, err := a.MessageService.FindSimplesAvailable(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 列表
func (a *MessageMController) GetManagedList(context *gin.Context) {
	messageTo := request.Param(context, "to").DefaultString("")
	messageType := request.Param(context, "type").DefaultString("")
	tenantId := request.Param(context, "tenantId").DefaultString("")
	messageScope := request.Param(context, "scope").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]interface{})
	if messageTo != "" {
		conds["to"] = messageTo
	}

	if messageType != "" {
		conds["type"] = messageTo
	}

	if messageScope != "" {
		conds["scope"] = messageScope
	}

	if search != "" {
		conds["search"] = search
	}

	if tenantId != "" {
		conds["tenantId"] = tenantId
	} else {
		conds["tenantId"] = "SYSTEM"
	}

	if count, lists, err := a.MessageService.FindSimples(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 新增
func (a *MessageMController) CreateMessage(c *gin.Context) {
	var data domain.Message
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if (data.ValidTime == nil || data.ValidTime.Before(time.Now())) &&
		(data.DueTime == nil || data.DueTime.After(time.Now())) {
		data.Status = 1
	}

	data.From = SecurityUtil.GetCurrentUserId(c)
	data.TenantId = "SYSTEM"

	if data.Scope == "" {
		data.Scope = "TENANT"
	}

	if result, err := a.MessageService.CreateMessage(&data); err == nil {
		response.Success(c, result)

		// 操作日志
		eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("CREATE").
			WithTarget("Message", result.Id).
			WithContent("${operator.name} 创建了消息【${target.type}】: ${target.title}"))
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// 修改
func (a *MessageMController) UpdateMessage(c *gin.Context) {
	var data domain.Message
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.MessageService.UpdateMessage(&data); err == nil {
		response.Success(c, result)

		// 操作日志
		eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("UPDATE").
			WithTarget("Message", result.Id).
			WithContent("${operator.name} 修改了消息【${target.type}】: ${target.title}"))
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 修改
func (a *MessageMController) PatchMessage(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var data = make(map[string]interface{})
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if result, err := a.MessageService.PatchMessage(id, data); err == nil {
		response.Success(c, result)

		// 操作日志
		eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("UPDATE").
			WithTarget("Message", result.Id).
			WithContent("${operator.name} 修改了消息【${target.type}】: ${target.title}"))
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 删除
func (a *MessageMController) DeleteMessage(c *gin.Context) {
	id, err := request.Param(c, "id").MustInt64()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	message, err := a.MessageService.GetById(id)
	if err != nil {
		response.SystemFail(c, err)
		return
	}
	if message == nil {
		response.Success(c, nil)
		return
	}

	if err := a.MessageService.DeleteMessage(message); err == nil {
		response.Success(c, "")
		// 操作日志
		eventbus.DispatchEvent("SYSTEM_LOG_OPERATION", domain.NewOperationLog("DELETE").
			WithTarget("Message", message.Id).
			WithContent("${operator.name} 删除了消息【${target.type}】: ${target.title}"))
	} else {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	}
}
