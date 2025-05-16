package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"
	"github.com/gophab/gophrame/service"
)

type TaskService struct {
	service.BaseService
	TaskRepository *repository.TaskRepository `inject:"taskRepository"`
}

var taskService = &TaskService{}

func init() {
	inject.InjectValue("taskService", taskService)
}

func (r *TaskService) GetById(id string) (*domain.Task, error) {
	return r.TaskRepository.GetById(id)
}

func (r *TaskService) FindByCreatedBy(createdBy string, pageable query.Pageable) (int64, []*domain.Task, error) {
	return r.TaskRepository.FindByCreatedBy(createdBy, pageable)
}

func (r *TaskService) CreateTask(task *domain.Task) (*domain.Task, error) {
	return r.TaskRepository.CreateTask(task)
}

func (r *TaskService) UpdateTask(task *domain.Task) (*domain.Task, error) {
	return r.TaskRepository.UpdateTask(task)
}

func (r *TaskService) DeleteTask(task *domain.Task) (*domain.Task, error) {
	return r.TaskRepository.DeleteTask(task)
}

func (r *TaskService) FinishTask(task *domain.Task) (*domain.Task, error) {

	if res := r.TaskRepository.Save(task); res.Error == nil && res.RowsAffected > 0 {
		return task, nil
	} else {
		return nil, res.Error
	}
}

func GetTask(taskId string) *domain.Task {
	task, _ := taskService.GetById(taskId)
	return task
}

func CreateTask(task *domain.Task) *domain.Task {
	result, err := taskService.CreateTask(task)
	if err == nil {
		return result
	}
	return nil
}

func UpdateTask(task *domain.Task) *domain.Task {
	result, err := taskService.UpdateTask(task)
	if err == nil {
		return result
	}
	return nil
}

func DeleteTask(task *domain.Task) *domain.Task {
	result, err := taskService.DeleteTask(task)
	if err == nil {
		return result
	}
	return nil
}
