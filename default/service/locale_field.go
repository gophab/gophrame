package service

import (
	"github.com/gophab/gophrame/core/i18n"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
	"github.com/gophab/gophrame/service"
)

type LocaleFieldService struct {
	service.BaseService
	LocaleFieldRepository *repository.LocaleFieldRepository `inject:"localeFieldRepository"`
	i18n.Translator
}

var localeFieldService = &LocaleFieldService{}

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
}

func (s *LocaleFieldService) LoadTranslations(entityName, entityId, locale string) []*i18n.LocaleFieldValue {
	entities := s.LocaleFieldRepository.GetLocaleAll(entityName, entityId, locale)
	if len(entities) > 0 {
		var results = make([]*i18n.LocaleFieldValue, len(entities))
		for i, fieldValue := range entities {
			results[i] = fieldValue.LocaleFieldValue
		}
		return results
	}
	return []*i18n.LocaleFieldValue{}
}

func (s *LocaleFieldService) SaveAll(fields []*domain.LocaleField) []*domain.LocaleField {
	return s.LocaleFieldRepository.SaveAll(fields)
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
