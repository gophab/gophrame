package openapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
	"github.com/gophab/gophrame/module/common/service"
)

type TaskOpenController struct {
	controller.ResourceController
	TaskService *service.TaskService `inject:"taskService"`
}

var taskOpenController = &TaskOpenController{}

type AdminTaskOpenController struct {
	controller.ResourceController
	TaskService *service.TaskService `inject:"taskService"`
}

var adminTaskOpenController = &AdminTaskOpenController{}

func init() {
	inject.InjectValue("taskOpenController", taskOpenController)
	inject.InjectValue("adminTaskOpenController", adminTaskOpenController)
}

func (m *TaskOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tasks", Handler: m.GetTasks},
		{HttpMethod: "GET", ResourcePath: "/task/:id", Handler: m.GetTask},
		{HttpMethod: "DELETE", ResourcePath: "/task/:id", Handler: m.DeleteTask},
	})
}

// GET /task/:id
func (c *TaskOpenController) GetTask(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.TaskService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result == nil {
		response.NotFound(ctx, "Not Found")
		return
	}

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if result.CreatedBy != currentUserId {
		response.NotAllowed(ctx, "Not Allowed")
		return
	}

	response.Success(ctx, result)
}

// DELETE /task/:id
func (c *TaskOpenController) DeleteTask(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.TaskService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result == nil {
		response.NotFound(ctx, "Not Found")
		return
	}

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if result.CreatedBy != currentUserId {
		response.NotAllowed(ctx, "Not Allowed")
		return
	}

	result, err = c.TaskService.DeleteTask(result)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	response.Success(ctx, result)
}

// GET /tasks
func (c *TaskOpenController) GetTasks(ctx *gin.Context) {
	pageable := query.GetPageable(ctx)

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if count, lists, err := c.TaskService.FindByCreatedBy(currentUserId, pageable); err == nil {
		ctx.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(ctx, lists)
	} else {
		ctx.Header("X-Total-Count", "0")
		response.Success(ctx, []any{})
	}
}

func (m *AdminTaskOpenController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tasks", Handler: m.GetTasks},
		{HttpMethod: "GET", ResourcePath: "/task/:id", Handler: m.GetTask},
		{HttpMethod: "DELETE", ResourcePath: "/task/:id", Handler: m.DeleteTask},
	})
}

// GET /task/:id
func (c *AdminTaskOpenController) GetTask(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.TaskService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result == nil {
		response.NotFound(ctx, "Not Found")
		return
	}

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if result.CreatedBy != currentUserId {
		response.NotAllowed(ctx, "Not Allowed")
		return
	}

	response.Success(ctx, result)
}

// DELETE /task/:id
func (c *AdminTaskOpenController) DeleteTask(ctx *gin.Context) {
	id, err := request.Param(ctx, "id").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.TaskService.GetById(id)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	if result == nil {
		response.NotFound(ctx, "Not Found")
		return
	}

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if result.CreatedBy != currentUserId {
		response.NotAllowed(ctx, "Not Allowed")
		return
	}

	result, err = c.TaskService.DeleteTask(result)
	if err != nil {
		response.SystemFail(ctx, err)
		return
	}

	response.Success(ctx, result)
}

// GET /tasks
func (c *AdminTaskOpenController) GetTasks(ctx *gin.Context) {
	pageable := query.GetPageable(ctx)

	currentUserId := SecurityUtil.GetCurrentUserId(ctx)

	if count, lists, err := c.TaskService.FindByCreatedBy(currentUserId, pageable); err == nil {
		ctx.Header("X-Total-Count", strconv.FormatInt(count, 10))
		response.Success(ctx, lists)
	} else {
		ctx.Header("X-Total-Count", "0")
		response.Success(ctx, []any{})
	}
}
