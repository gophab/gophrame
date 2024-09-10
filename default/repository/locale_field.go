package repository

import (
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/default/domain"
	"gorm.io/gorm"
)

type LocaleFieldRepository struct {
	*gorm.DB `inject:"database"`
}

var localeFieldRepository = &LocaleFieldRepository{}

func init() {
	inject.InjectValue("localeFieldRepository", localeFieldRepository)
}

func (r *LocaleFieldRepository) SaveAll(entities []*domain.LocaleField) []*domain.LocaleField {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 如果有异常,执行回滚事务
		}
	}()

	for _, entity := range entities {
		tx.Save(entity)
	}

	tx.Commit()

	return entities
	// if res := r.DB.CreateInBatches(entities, 50); res.Error == nil && res.RowsAffected > 0 {
	// 	return entities
	// } else {
	// 	return nil
	// }
}

func (r *LocaleFieldRepository) GetLocaleAll(locale string, entityName string, entityIds ...string) []*domain.LocaleField {
	tx := r.DB.Model(&domain.LocaleField{}).
		Where("entity_name = ?", entityName)

	if len(entityIds) > 1 {
		tx = tx.Where("entity_id in ?", entityIds)
	} else {
		tx = tx.Where("entity_id = ?", entityIds[0])
	}

	tx = tx.Where("locale = ?", locale)

	var results []*domain.LocaleField
	if res := tx.Find(&results); res.Error == nil {
		return results
	} else {
		return nil
	}
}

func (r *LocaleFieldRepository) GetEntityAll(entityName string, entityId string) []*domain.LocaleField {
	tx := r.DB.Model(&domain.LocaleField{}).
		Where("entity_name = ?", entityName)

	ids := strings.Split(entityId, ",")
	if len(ids) > 1 {
		tx = tx.Where("entity_id in ?", ids)
	} else {
		tx = tx.Where("entity_id = ?", entityId)
	}

	var results []*domain.LocaleField
	if res := tx.Find(&results); res.Error == nil {
		return results
	} else {
		return nil
	}
}

func (r *LocaleFieldRepository) GetEntityFieldAll(entityName string, entityId string, field string) []*domain.LocaleField {
	tx := r.DB.Model(&domain.LocaleField{}).
		Where("entity_name = ?", entityName)

	ids := strings.Split(entityId, ",")
	if len(ids) > 1 {
		tx = tx.Where("entity_id in ?", ids)
	} else {
		tx = tx.Where("entity_id = ?", entityId)
	}

	tx = tx.Where("name = ?", field)

	var results []*domain.LocaleField
	if res := tx.Find(&results); res.Error == nil {
		return results
	} else {
		return nil
	}
}
