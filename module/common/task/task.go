package task

import (
	"time"

	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/service"
)

func CreateAutoTask(createdBy string, name string, description string) string {
	var task *domain.Task = &domain.Task{
		CreatedBy:   createdBy,
		Name:        name,
		Description: util.StringAddr(description),
		Type:        "AUTO",
		Status:      0,
		CreatedTime: time.Now(),
	}

	result := service.CreateTask(task)
	if result != nil {
		return result.Id
	}
	return ""
}

func CreateManualTask(createdBy string, name string, description string) string {
	var task *domain.Task = &domain.Task{
		CreatedBy:   createdBy,
		Name:        name,
		Description: util.StringAddr(description),
		Type:        "MANUAL",
		Status:      0,
		CreatedTime: time.Now(),
	}

	result := service.CreateTask(task)
	if result != nil {
		return result.Id
	}
	return ""
}

func CreateAsyncTask(createdBy string, name string, description string) string {
	var task *domain.Task = &domain.Task{
		CreatedBy:   createdBy,
		Name:        name,
		Description: util.StringAddr(description),
		Type:        "ASYNC",
		Status:      0,
		CreatedTime: time.Now(),
	}

	result := service.CreateTask(task)
	if result != nil {
		return result.Id
	}
	return ""
}

func StartTask(taskId string) {
	task := service.GetTask(taskId)
	if task != nil {
		task.Status = 1 /* STARTED */
		service.UpdateTask(task)
	}
}

func UpdateTask(taskId string, progress float32, status int, remark string) {
	task := service.GetTask(taskId)
	if task != nil {
		if status > 0 {
			task.Status = status /* FINISHED */
		}

		if progress > 0 {
			task.Progress = progress /* 进度 */
		}

		if remark != "" {
			task.Remark = util.StringAddr(remark)
		}

		service.UpdateTask(task)
	}
}

func FinishTask(taskId string, mode string, result string) {
	task := service.GetTask(taskId)
	if task != nil {
		task.Mode = mode
		task.Result = util.StringAddr(result)
		task.Progress = 100 /* 进度 */
		task.Status = 2     /* FINISHED */
		task.FinishedTime = util.TimeAddr(time.Now())
		service.UpdateTask(task)
	}
}
