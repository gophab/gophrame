package repository

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/module/common/domain"
	"gorm.io/gorm"
)

type TaskRepository struct {
	*gorm.DB `inject:"database"`
}

var taskRepository = &TaskRepository{}

func init() {
	inject.InjectValue("taskRepository", taskRepository)
}

func (r *TaskRepository) GetById(id string) (*domain.Task, error) {
	var result domain.Task
	if res := r.Model(&domain.Task{}).Where("id=?", id).Where("del_flag = ?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *TaskRepository) FindByCreatedBy(createdBy string, pageable query.Pageable) (int64, []*domain.Task, error) {
	var results = make([]*domain.Task, 0)
	var q = r.Model(&domain.Task{}).Where("created_by=?", createdBy).Where("del_flag = ?", false)
	var count int64 = 0
	if !pageable.NoCount() {
		if res := q.Count(&count); res.Error != nil {
			return 0, nil, res.Error
		}
	}

	if res := query.Page(q, pageable).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return count, results, nil
	} else {
		return 0, nil, res.Error
	}
}

func (r *TaskRepository) CreateTask(task *domain.Task) (*domain.Task, error) {
	if res := r.Save(task); res.Error == nil && res.RowsAffected > 0 {
		return task, nil
	} else {
		return nil, res.Error
	}
}

func (r *TaskRepository) UpdateTask(task *domain.Task) (*domain.Task, error) {
	if res := r.Save(task); res.Error == nil && res.RowsAffected > 0 {
		return task, nil
	} else {
		return nil, res.Error
	}
}

func (r *TaskRepository) DeleteTask(task *domain.Task) (*domain.Task, error) {
	task.DelFlag = true
	if res := r.Save(task); res.Error == nil && res.RowsAffected > 0 {
		return task, nil
	} else {
		return nil, res.Error
	}
}
