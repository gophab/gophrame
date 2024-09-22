package openapi

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

type MessageOpenController struct {
	controller.ResourceController
	MessageService *service.MessageService `inject:"messageService"`
}

var messageOpenController *MessageOpenController = &MessageOpenController{}

type AdminMessageOpenController struct {
	controller.ResourceController
	MessageService *service.MessageService `inject:"messageService"`
}

var adminMessageOpenController *AdminMessageOpenController = &AdminMessageOpenController{}

func init() {
	inject.InjectValue("messageOpenController", messageOpenController)
	inject.InjectValue("adminMessageOpenController", adminMessageOpenController)
}

// 组织
func (m *MessageOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/messages", Handler: m.GetList},
		{HttpMethod: "GET", ResourcePath: "/messages/managed", Handler: m.GetManagedList},
		{HttpMethod: "GET", ResourcePath: "/message/:id", Handler: m.GetMessage},
		{HttpMethod: "POST", ResourcePath: "/message", Handler: m.CreateMessage},
		{HttpMethod: "PUT", ResourcePath: "/message", Handler: m.UpdateMessage},
		{HttpMethod: "DELETE", ResourcePath: "/message/:id", Handler: m.DeleteMessage},
	})
}

// 列表
func (a *MessageOpenController) GetList(context *gin.Context) {
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
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.MessageService.FindSimplesAvailable(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 列表: 我可以管理的
func (a *MessageOpenController) GetManagedList(context *gin.Context) {
	messageTo := request.Param(context, "to").DefaultString("")
	messageType := request.Param(context, "type").DefaultString("")
	search := request.Param(context, "search").DefaultString("")
	pageable := query.GetPageable(context)

	var conds = make(map[string]interface{})
	if messageTo != "" {
		conds["to"] = messageTo
	}

	if messageType != "" {
		conds["type"] = messageType
	}

	if search != "" {
		conds["search"] = search
	}

	conds["scope"] = "PRIVATE"
	conds["from"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.MessageService.FindSimples(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 1.根据id查询节点
func (a *MessageOpenController) GetMessage(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	show := request.Param(context, "show").DefaultBool(true)
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	result, err := a.MessageService.GetById(id)
	if err != nil {
		response.SystemFail(context, err)
		return
	}

	if result == nil {
		response.NotFound(context, "Not Found")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(context)
	tenantId := SecurityUtil.GetCurrentTenantId(context)

	switch result.Scope {
	case "PRIVATE":
		if result.To != userId && result.From != userId {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	case "TENANT":
		if result.TenantId != tenantId {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	case "PUBLIC":
		if result.TenantId != tenantId && result.TenantId != "SYSTEM" {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	}

	response.Success(context, result)

	if show {
		eventbus.DispatchEvent("SYSTEM_MESSAGE_VIEWED", result, userId)
	}
}

// 新增
func (a *MessageOpenController) CreateMessage(c *gin.Context) {
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
	data.Scope = "PRIVATE"
	data.TenantId = SecurityUtil.GetCurrentTenantId(c)

	if result, err := a.MessageService.CreateMessage(&data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// 修改
func (a *MessageOpenController) UpdateMessage(c *gin.Context) {
	var data domain.Message
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := a.MessageService.GetById(data.Id)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	if result == nil {
		response.NotFound(c, "Not Found")
		return
	}

	if result.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(c)

	switch result.Scope {
	case "PRIVATE":
		if result.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	default:
		response.NotAllowed(c, "Not Allowed")
		return
	}

	result.Title = data.Title
	result.Content = data.Content
	result.ValidTime = data.ValidTime
	result.DueTime = data.DueTime

	if result, err := a.MessageService.UpdateMessage(result); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 修改
func (a *MessageOpenController) PatchMessage(c *gin.Context) {
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
	if message.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	var data = make(map[string]interface{})
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := a.MessageService.GetById(id)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	if result == nil {
		response.NotFound(c, "Not Found")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(c)

	switch result.Scope {
	case "PRIVATE":
		if result.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	default:
		response.NotAllowed(c, "Not Allowed")
		return
	}

	// 禁止修改字段
	delete(data, "from")
	delete(data, "scope")
	delete(data, "tenantId")
	delete(data, "tenant_id")

	if result, err := a.MessageService.PatchMessage(id, data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 删除
func (a *MessageOpenController) DeleteMessage(c *gin.Context) {
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

	userId := SecurityUtil.GetCurrentUserId(c)

	switch message.Scope {
	case "PRIVATE":
		if message.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	default:
		response.NotAllowed(c, "Not Allowed")
		return
	}

	if err := a.MessageService.DeleteMessage(message); err == nil {
		response.Success(c, "")
	} else {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	}
}

// 组织
func (m *AdminMessageOpenController) AfterInitialize() {
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

// 列表: Available
func (a *AdminMessageOpenController) GetList(context *gin.Context) {
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
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.MessageService.FindSimplesAvailable(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 列表
func (a *AdminMessageOpenController) GetManagedList(context *gin.Context) {
	messageTo := request.Param(context, "to").DefaultString("")
	messageType := request.Param(context, "type").DefaultString("")
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

	conds["userId"] = SecurityUtil.GetCurrentUserId(context)
	conds["tenantId"] = SecurityUtil.GetCurrentTenantId(context)

	if count, lists, err := a.MessageService.FindSimples(conds, pageable); err == nil {
		context.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(context, lists)
	} else {
		context.Header("X-Total-Count", "0")
		response.Success(context, []any{})
	}
}

// 1.根据id查询节点
func (a *AdminMessageOpenController) GetMessage(context *gin.Context) {
	id, err := request.Param(context, "id").MustInt64()
	show := request.Param(context, "show").DefaultBool(true)
	if err != nil {
		response.FailCode(context, errors.INVALID_PARAMS)
		return
	}

	result, err := a.MessageService.GetById(id)
	if err != nil {
		response.SystemFail(context, err)
		return
	}

	if result == nil {
		response.NotFound(context, "Not Found")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(context)
	tenantId := SecurityUtil.GetCurrentTenantId(context)

	switch result.Scope {
	case "PRIVATE":
		if result.To != userId && result.From != userId {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	case "TENANT":
		if result.TenantId != tenantId {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	case "PUBLIC":
		if result.TenantId != tenantId && result.TenantId != "SYSTEM" {
			response.NotAllowed(context, "Not Allowed")
			return
		}
	}

	response.Success(context, result)

	if show {
		eventbus.DispatchEvent("SYSTEM_MESSAGE_VIEWED", result, userId)
	}
}

// 新增
func (a *AdminMessageOpenController) CreateMessage(c *gin.Context) {
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
	data.TenantId = SecurityUtil.GetCurrentTenantId(c)

	if data.Scope == "" {
		data.Scope = "TENANT"
	}

	if data.Scope == "PUBLIC" {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	if result, err := a.MessageService.CreateMessage(&data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
	}
}

// 修改
func (a *AdminMessageOpenController) UpdateMessage(c *gin.Context) {
	var data domain.Message
	if err := c.ShouldBind(&data); err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	result, err := a.MessageService.GetById(data.Id)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	if result == nil {
		response.NotFound(c, "Not Found")
		return
	}

	if result.TenantId != SecurityUtil.GetCurrentTenantId(c) {
		response.NotAllowed(c, "Not Allowed")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(c)
	tenantId := SecurityUtil.GetCurrentTenantId(c)

	switch result.Scope {
	case "PRIVATE":
		if result.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	case "TENANT", "PUBLIC":
		if result.TenantId != tenantId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	}

	result.Title = data.Title
	result.Content = data.Content
	result.ValidTime = data.ValidTime
	result.DueTime = data.DueTime

	if result, err := a.MessageService.UpdateMessage(result); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 修改
func (a *AdminMessageOpenController) PatchMessage(c *gin.Context) {
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

	result, err := a.MessageService.GetById(id)
	if err != nil {
		response.SystemFail(c, err)
		return
	}

	if result == nil {
		response.NotFound(c, "Not Found")
		return
	}

	userId := SecurityUtil.GetCurrentUserId(c)
	tenantId := SecurityUtil.GetCurrentTenantId(c)

	switch result.Scope {
	case "PRIVATE":
		if result.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	case "TENANT", "PUBLIC":
		if result.TenantId != tenantId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	}

	// 禁止修改字段
	delete(data, "from")
	delete(data, "scope")
	delete(data, "tenantId")
	delete(data, "tenant_id")

	if result, err := a.MessageService.PatchMessage(id, data); err == nil {
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
	}
}

// 删除
func (a *AdminMessageOpenController) DeleteMessage(c *gin.Context) {
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

	userId := SecurityUtil.GetCurrentUserId(c)
	tenantId := SecurityUtil.GetCurrentTenantId(c)

	switch message.Scope {
	case "PRIVATE":
		if message.From != userId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	case "TENANT", "PUBLIC":
		if message.TenantId != tenantId {
			response.NotAllowed(c, "Not Allowed")
			return
		}
	}

	if err := a.MessageService.DeleteMessage(message); err == nil {
		response.Success(c, "")
	} else {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	}
}
