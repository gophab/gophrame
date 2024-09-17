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

type MessageRepository struct {
	*gorm.DB `inject:"database"`
}

var messageRepository = &MessageRepository{}

func init() {
	inject.InjectValue("messageRepository", messageRepository)
}

func (r *MessageRepository) CreateMessage(message *domain.Message) (*domain.Message, error) {
	if message.ValidTime == nil {
		message.ValidTime = util.TimeAddr(time.Now())
	}

	if res := r.Model(&domain.Message{}).Create(message); res.Error == nil {
		return message, nil
	} else {
		return nil, res.Error
	}
}

func (r *MessageRepository) GetById(id int64) (*domain.Message, error) {
	var result domain.Message
	if res := r.Model(&domain.Message{}).Where("id = ?", id).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *MessageRepository) PatchMessage(id int64, data map[string]interface{}) (*domain.Message, error) {
	data["id"] = id
	if res := r.Model(&domain.Message{}).Where("id=?", id).UpdateColumns(data); res.Error == nil {
		return r.GetById(id)
	} else {
		return nil, res.Error
	}
}

func (r *MessageRepository) DeleteMessage(message *domain.Message) error {
	return r.Model(&domain.Message{}).Where("id = ?", message.Id).Update("del_flag", true).Error
}

func (r *MessageRepository) DeleteById(id int64) error {
	return r.Model(&domain.Message{}).Where("id = ?", id).Update("del_flag", true).Error
}

func (r *MessageRepository) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Message, error) {
	tx := r.Model(&domain.Message{})
	for k, v := range conds {
		if k == "user_id" {
			// managed by user_id
			tx.Where("(`from` = ? and scope = 'PRIVATE') or (scope <> 'PRIVATE')", v)
		} else if k == "search" {
			tx.Where("title like ? or content like ?", "%"+v.(string)+"%", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}

	var count int64
	if !pageable.NoCount() {
		tx.Count(&count)
	}

	var list = make([]*domain.Message, 0)
	if res := query.Page(tx, pageable).Find(&list); res.Error == nil {
		return count, list, nil
	} else {
		return 0, []*domain.Message{}, res.Error
	}
}

func (r *MessageRepository) FindAvailable(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Message, error) {
	var userId = conds["user_id"]
	var tenantId = conds["tenant_id"]
	delete(conds, "user_id")
	delete(conds, "tenant_id")

	tx := r.Model(&domain.Message{})
	if userId != nil && userId.(string) != "" {
		if tenantId != nil && tenantId.(string) != "" {
			// user available: 包括私信、企业内部、公共
			tx.Where(
				tx.Where("(`to` = ? and scope = 'PRIVATE')", userId).
					Or("tenant_id = 'SYSTEM' and scope = 'PUBLIC'").
					Or("tenant_id = ? and scope = 'TENANT'"))
		} else {
			// 只有私信、公共
			tx.Where(
				tx.Where("(`to` = ? and scope = 'PRIVATE')", userId).
					Or("tenant_id = 'SYSTEM' and scope = 'PUBLIC'"))
		}
	} else if tenantId != nil && tenantId.(string) != "" {
		// 包括公共和企业内部
		tx.Where(
			tx.Where("tenant_id = 'SYSTEM' and scope = 'PUBLIC'").
				Or("tenant_id = ? and scope = 'TENANT'"))
	}

	for k, v := range conds {
		if k == "search" {
			tx.Where("title like ? or content like ?", "%"+v.(string)+"%", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}

	tx.Where("valid_time <= ?", time.Now()).
		Where("due_time is null or due_time >= ?", time.Now()).
		Where("status = ?", 1)

	var count int64
	if !pageable.NoCount() {
		//
		tx.Count(&count)
	}

	if userId != "" {
		tx.Select("*, EXISTS (?) AS `read`",
			gorm.Expr("select 1 from sys_message_access_log b where user_id=? and action=? and b.message_id=sys_message.id", userId, "READ"))
	}

	var list = make([]*domain.Message, 0)
	if res := query.Page(tx, pageable).Find(&list); res.Error == nil {
		return count, list, nil
	} else {
		return 0, []*domain.Message{}, res.Error
	}
}

func (r *MessageRepository) FindSimples(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.SimpleMessage, error) {
	tx := r.Model(&domain.Message{})
	for k, v := range conds {
		if k == "user_id" {
			// managed by user_id
			tx.Where("(`from` = ? and scope = 'PRIVATE') or (scope <> 'PRIVATE')", v)
		} else if k == "search" {
			tx.Where("title like ? or content like ?", "%"+v.(string)+"%", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}

	var count int64
	if !pageable.NoCount() {
		tx.Count(&count)
	}

	var list = make([]*domain.SimpleMessage, 0)
	if res := query.Page(tx, pageable).Find(&list); res.Error == nil {
		return count, list, nil
	} else {
		return 0, []*domain.SimpleMessage{}, res.Error
	}
}

func (r *MessageRepository) FindSimplesAvailable(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.SimpleMessage, error) {
	tx := r.Model(&domain.Message{})

	var userId = conds["user_id"]
	var tenantId = conds["tenant_id"]
	delete(conds, "user_id")
	delete(conds, "tenant_id")

	if userId != nil && userId.(string) != "" {
		if tenantId != nil && tenantId.(string) != "" {
			// user available: 包括私信、企业内部、公共
			tx.Where(
				tx.Where("(`to` = ? and scope = 'PRIVATE')", userId).
					Or("tenant_id = 'SYSTEM' and scope = 'PUBLIC'").
					Or("tenant_id = ? and scope = 'TENANT'", tenantId))
		} else {
			// 只有私信、公共
			tx.Where(
				tx.Where("(`to` = ? and scope = 'PRIVATE')", userId).
					Or("tenant_id = 'SYSTEM' and scope = 'PUBLIC'"))
		}
	} else if tenantId != nil && tenantId.(string) != "" {
		// 包括公共和企业内部
		tx.Where(
			tx.Where("tenant_id = 'SYSTEM' and scope = 'PUBLIC'").
				Or("tenant_id = ? and scope = 'TENANT'", tenantId))
	}

	for k, v := range conds {
		if k == "search" {
			tx.Where("title like ? or content like ?", "%"+v.(string)+"%", "%"+v.(string)+"%")
		} else {
			tx.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}

	tx.Where("valid_time is null or valid_time <= ?", time.Now()).
		Where("due_time is null or due_time >= ?", time.Now()).
		Where("status = ?", 1)

	var count int64
	if !pageable.NoCount() {
		//
		tx.Count(&count)
	}

	if userId != "" {
		tx.Select("*, EXISTS (?) AS `read`",
			gorm.Expr("select 1 from sys_message_access_log b where user_id=? and action=? and b.message_id=sys_message.id", userId, "READ"))
	}

	var list = make([]*domain.SimpleMessage, 0)
	if res := query.Page(tx, pageable).Find(&list); res.Error == nil {
		return count, list, nil
	} else {
		return 0, []*domain.SimpleMessage{}, res.Error
	}
}

func (r *MessageRepository) ValidateMessages() {
	// valid
	tx := r.Model(&domain.Message{}).
		Where("status = ?", 0).
		Where("valid_time <= ? or valid_time is null", time.Now()).
		Where("due_time > ? or due_time is null", time.Now())

	pageable := &query.Pagination{
		Page: 1,
		Size: 100,
		Sort: []query.Sort{
			{By: "id", Direction: "asc"},
		},
	}

	for {
		var list = make([]*domain.Message, 0)
		if res := query.Page(tx, pageable).Find(&list); res.Error == nil && res.RowsAffected > 0 {
			for _, message := range list {
				message.Status = 1
			}
			r.Updates(list)

			pageable.Page++
		} else {
			break
		}
	}

	// expired
	tx = r.Model(&domain.Message{}).
		Where("status = ?", 1).
		Where("due_time <= ?", time.Now())

	pageable = &query.Pagination{
		Page: 1,
		Size: 100,
		Sort: []query.Sort{
			{By: "id", Direction: "asc"},
		},
	}

	for {
		var list = make([]*domain.Message, 0)
		if res := query.Page(tx, pageable).Find(&list); res.Error == nil && res.RowsAffected > 0 {
			for _, message := range list {
				message.Status = -1
			}
			r.Updates(list)

			pageable.Page++
		} else {
			break
		}
	}
}

func (r *MessageRepository) HistoryMessages() {
	// valid
	tx := r.Model(&domain.Message{}).
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
		var list = make([]*domain.Message, 0)
		if res := query.Page(tx, pageable).Find(&list); res.Error == nil && res.RowsAffected > 0 {
			var histories = make([]*domain.MessageHistory, len(list))
			for index, message := range list {
				var history = domain.MessageHistory{Message: *message}
				histories[index] = &history

				message.Status = -2
			}
			r.Model(&domain.MessageHistory{}).CreateInBatches(histories, pageable.Size)
			r.Model(&domain.Message{}).Updates(list)

			pageable.Page++
		} else {
			break
		}
	}

	r.Delete(&domain.Message{}, "status = ?", -2)
}
