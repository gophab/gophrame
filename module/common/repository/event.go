package repository

import (
	"fmt"
	"time"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/common/domain"

	"gorm.io/gorm"
)

type EventRepository struct {
	*gorm.DB `inject:"database"`
}

var eventRepository = &EventRepository{}

func init() {
	inject.InjectValue("eventRepository", eventRepository)
}

func (r *EventRepository) CreateEvent(event *domain.Event) (*domain.Event, error) {
	if res := r.Model(&domain.Event{}).Create(event); res.Error == nil {
		return event, nil
	} else {
		return nil, res.Error
	}
}

func (r *EventRepository) GetById(id int64) (*domain.Event, error) {
	var result domain.Event
	if res := r.Model(&domain.Event{}).Where("id = ?", id).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *EventRepository) PatchEvent(id int64, data map[string]any) (*domain.Event, error) {
	data["id"] = id
	if res := r.Model(&domain.Event{}).Where("id=?", id).UpdateColumns(util.DbFields(data)); res.Error == nil {
		return r.GetById(id)
	} else {
		return nil, res.Error
	}
}

func (r *EventRepository) DeleteEvent(event *domain.Event) error {
	return r.Model(&domain.Event{}).Where("id = ?", event.Id).Update("del_flag", true).Error
}

func (r *EventRepository) DeleteById(id int64) error {
	return r.Model(&domain.Event{}).Where("id = ?", id).Update("del_flag", true).Error
}

func (r *EventRepository) Find(conds map[string]any, pageable query.Pageable) (int64, []*domain.Event, error) {
	tx := r.Model(&domain.Event{})

	var userId = conds["user_id"]
	delete(conds, "user_id")
	for k, v := range conds {
		if k == "search" {
			tx.Where("content like ?", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("%s = ?", k), v)
		}
	}

	var count int64
	if !pageable.NoCount() {
		tx.Count(&count)
	}

	tx.Order("created_time desc")

	var list = make([]*domain.Event, 0)
	if res := query.Page(tx, pageable).Find(&list); res.Error == nil {
		if userId != "" {
			var accessTime = time.Now()
			if res := r.Table("sys_event_access_log b").
				Select("access_time").
				Where("user_id=?", userId).
				Where("action=?", "READ").
				First(&accessTime); res.Error == nil && res.RowsAffected > 0 {
				for _, event := range list {
					event.Read = event.CreatedTime.Before(accessTime)
				}
			}
		}
		return count, list, nil
	} else {
		return 0, []*domain.Event{}, res.Error
	}
}

func (r *EventRepository) Check(conds map[string]any) (*domain.Event, error) {
	tx := r.Model(&domain.Event{})

	var userId = conds["user_id"]
	delete(conds, "user_id")

	for k, v := range conds {
		if k == "search" {
			tx.Where("content like ?", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("%s = ?", k), v)
		}
	}
	tx.Order("created_time desc")

	var result domain.Event
	if res := tx.First(&result); res.Error == nil && res.RowsAffected > 0 {
		if userId != "" {
			var accessTime = time.Now()
			if res := r.Table("sys_event_access_log b").
				Select("access_time").
				Where("user_id=?", userId).
				Where("action=?", "READ").
				First(&accessTime); res.Error == nil && res.RowsAffected > 0 {
				result.Read = result.CreatedTime.Before(accessTime)
			}
		}
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *EventRepository) HistoryEvents() {
	// valid
	tx := r.Model(&domain.Event{}).
		Where("status = ?", -1).
		Or("created_time < ?", time.Now().Add(-time.Hour*24*180)).
		Or("del_flag = ?", true)

	pageable := &query.Pagination{
		Page: 1,
		Size: 100,
		Sort: []query.Sort{
			{By: "id", Direction: "asc"},
		},
	}

	for {
		var list = make([]*domain.Event, 0)
		if res := query.Page(tx, pageable).Find(&list); res.Error == nil && res.RowsAffected > 0 {
			var histories = make([]*domain.EventHistory, len(list))
			for index, event := range list {
				var history = domain.EventHistory{Event: *event}
				histories[index] = &history

				event.Status = -2
			}
			r.Model(&domain.EventHistory{}).CreateInBatches(histories, pageable.Size)
			r.Model(&domain.Event{}).Updates(list)

			pageable.Page++
		} else {
			break
		}
	}

	r.Delete(&domain.Event{}, "status = ?", -2)
}
