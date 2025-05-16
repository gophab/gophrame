package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"
)

type SysOptionService struct {
	service.BaseService
	SysOptionRepository *repository.SysOptionRepository `inject:"sysOptionRepository"`
}

var sysOptionService = &SysOptionService{}

func init() {
	inject.InjectValue("sysOptionService", sysOptionService)
}

var defaultSystemOptions = map[string]string{}

func (s *SysOptionService) GetDefaultOptions(tenantId string) (*domain.SysOptions, error) {
	result := &domain.SysOptions{TenantId: tenantId, Options: make(map[string]*domain.SysOption)}

	if tenantId != "DEFAULT" {
		// DEFAULT by code
		for k, v := range defaultSystemOptions {
			result.Options[k] = &domain.SysOption{
				TenantId: tenantId,
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}
		}

		// DEFAULT in DB
		resultDB, err := s.SysOptionRepository.GetDefaultOptions()
		if err != nil {
			return nil, err
		}

		if resultDB != nil && len(resultDB.Options) > 0 {
			for k, v := range resultDB.Options {
				v.TenantId = tenantId
				result.Options[k] = v
			}
			return result, nil
		}
	}

	return result, nil
}

func (s *SysOptionService) GetTenantOptions(tenantId string) (*domain.SysOptions, error) {
	result, err := s.GetDefaultOptions(tenantId)
	if err != nil {
		return nil, err
	}

	// Tenant Options in DB
	resultDB, err := s.SysOptionRepository.GetTenantOptions(tenantId)
	if err != nil {
		return nil, err
	}

	result.TenantId = tenantId

	if resultDB != nil && len(resultDB.Options) > 0 {
		for k, v := range resultDB.Options {
			v.TenantId = tenantId
			result.Options[k] = v
		}
	}

	return result, nil
}

func (s *SysOptionService) GetSystemOptions() (*domain.SysOptions, error) {
	return s.GetTenantOptions("SYSTEM")
}

func (s *SysOptionService) AddSysOption(option *domain.SysOption) (*domain.SysOption, error) {
	if res := s.SysOptionRepository.Save(option); res.Error == nil && res.RowsAffected > 0 {
		return option, nil
	} else {
		return nil, res.Error
	}
}

func (s *SysOptionService) DeleteSysOption(option *domain.SysOption) (*domain.SysOption, error) {
	if res := s.SysOptionRepository.Delete(&domain.SysOption{}).Where("name = ? and tenant_id = ?", option.Name, option.TenantId); res.Error == nil {
		return option, nil
	} else {
		return nil, res.Error
	}
}

func (s *SysOptionService) AddSysOptions(options []*domain.SysOption) ([]*domain.SysOption, error) {
	var result = make([]*domain.SysOption, len(options))
	for i, option := range options {
		if res := s.SysOptionRepository.Save(option); res.Error != nil {
			//if res := s.SysOptionRepository.Updates(option); res.Error != nil {
			return nil, res.Error
			//}
		}
		result[i] = option
	}
	return result, nil
}

func (s *SysOptionService) RemoveAllTenantOptions(tenantId string) error {
	return s.SysOptionRepository.RemoveAllTenantOptions(tenantId)
}

func (s *SysOptionService) RemoveTenantOption(tenantId string, key string) (*domain.SysOption, error) {
	return nil, s.SysOptionRepository.Delete(&domain.SysOption{}, "tenant_id = ? and name = ?", tenantId, key).Error
}

func (s *SysOptionService) SetTenantOption(tenantId string, key string, value string) (*domain.SysOption, error) {
	var option = domain.SysOption{
		TenantId: tenantId,
		Option: domain.Option{
			Name:      key,
			Value:     value,
			ValueType: "STRING",
		},
	}

	if res := s.SysOptionRepository.Save(&option); res.Error == nil && res.RowsAffected > 0 {
		return &option, nil
	} else {
		return nil, res.Error
	}
}

func (s *SysOptionService) SetTenantOptions(tenantOptions *domain.SysOptions) (*domain.SysOptions, error) {
	// 1. Remove Sys Options
	if err := s.RemoveAllTenantOptions(tenantOptions.TenantId); err != nil {
		return nil, err
	}

	// 2. Save
	var options = make([]*domain.SysOption, 0)
	for _, v := range tenantOptions.Options {

		v.TenantId = tenantOptions.TenantId
		options = append(options, v)
	}

	if _, err := s.AddSysOptions(options); err != nil {
		return nil, err
	}

	return tenantOptions, nil
}

func (s *SysOptionService) UpdateTenantOptions(tenantOptions *domain.SysOptions) (*domain.SysOptions, error) {
	// 2. Save
	var options = make([]*domain.SysOption, 0)
	for _, v := range tenantOptions.Options {
		v.TenantId = tenantOptions.TenantId
		// 1. Remove Sys Options
		s.RemoveTenantOption(v.TenantId, v.Name)
		options = append(options, v)
	}

	if _, err := s.AddSysOptions(options); err != nil {
		return nil, err
	}

	return tenantOptions, nil
}
