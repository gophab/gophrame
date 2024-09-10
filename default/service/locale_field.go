package service

import (
	"time"

	"github.com/gophab/gophrame/core/i18n"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
	"github.com/gophab/gophrame/service"

	"github.com/patrickmn/go-cache"
)

type LocaleFieldService struct {
	service.BaseService
	i18n.Translator

	LocaleFieldRepository *repository.LocaleFieldRepository `inject:"localeFieldRepository"`
	LocalCache            *cache.Cache
}

var localeFieldService = &LocaleFieldService{
	LocalCache: cache.New(30*time.Minute, 1*time.Minute),
}

func init() {
	inject.InjectValue("localeFieldService", localeFieldService)
	inject.InjectValue("translator", localeFieldService)
}

func (s *LocaleFieldService) StoreTranslations(translations []*i18n.LocaleFieldValue) {
	var fields = make([]*domain.LocaleField, len(translations))
	for i, fieldValue := range translations {
		fields[i] = &domain.LocaleField{
			LocaleFieldValue: fieldValue,
		}
	}
	s.LocaleFieldRepository.SaveAll(fields)

	for _, field := range fields {
		var key = field.EntityName + ":" + field.EntityId + ":" + field.Locale
		s.LocalCache.Delete(key)
	}
}

func (s *LocaleFieldService) LoadTranslations(locale, entityName string, entityIds ...string) map[string][]*i18n.LocaleFieldValue {
	var results = make(map[string][]*i18n.LocaleFieldValue)

	var loadIds = make([]string, 0)
	for _, entityId := range entityIds {
		var key = entityName + ":" + entityId + ":" + locale
		if v, b := s.LocalCache.Get(key); b && v != nil {
			results[entityId] = v.([]*i18n.LocaleFieldValue)
		} else {
			loadIds = append(loadIds, entityId)
		}
	}

	if len(loadIds) > 0 {
		entityFields := s.LocaleFieldRepository.GetLocaleAll(locale, entityName, loadIds...)
		if len(entityFields) > 0 {
			for _, fieldValue := range entityFields {
				list, b := results[fieldValue.EntityId]
				if !b {
					list = make([]*i18n.LocaleFieldValue, 0)
					list = append(list, fieldValue.LocaleFieldValue)
				}
				results[fieldValue.EntityId] = list

				var key = entityName + ":" + fieldValue.EntityId + ":" + locale
				s.LocalCache.Set(key, list, cache.DefaultExpiration)
			}
		}

		// 对于没有I18n信息的实体，缓存为空，避免再次查询
		for _, entityId := range loadIds {
			if _, b := results[entityId]; !b {
				var key = entityName + ":" + entityId + ":" + locale
				s.LocalCache.Set(key, []*i18n.LocaleFieldValue{}, cache.DefaultExpiration)
			}
		}
	}
	return results
}

func (s *LocaleFieldService) SaveAll(fields []*domain.LocaleField) []*domain.LocaleField {
	var results = s.LocaleFieldRepository.SaveAll(fields)
	for _, field := range fields {
		var key = field.EntityName + ":" + field.EntityId + ":" + field.Locale
		s.LocalCache.Delete(key)
	}
	return results
}

func (s *LocaleFieldService) GetEntityAll(entityName, entityId string) []*i18n.LocaleFieldValue {
	entities := s.LocaleFieldRepository.GetEntityAll(entityName, entityId)
	if len(entities) > 0 {
		var results = make([]*i18n.LocaleFieldValue, len(entities))
		for i, fieldValue := range entities {
			results[i] = fieldValue.LocaleFieldValue
		}
		return results
	}
	return []*i18n.LocaleFieldValue{}
}

func (s *LocaleFieldService) GetEntityFieldAll(entityName, entityId, field string) []*i18n.LocaleFieldValue {
	entities := s.LocaleFieldRepository.GetEntityFieldAll(entityName, entityId, field)
	if len(entities) > 0 {
		var results = make([]*i18n.LocaleFieldValue, len(entities))
		for i, fieldValue := range entities {
			results[i] = fieldValue.LocaleFieldValue
		}
		return results
	}
	return []*i18n.LocaleFieldValue{}
}
