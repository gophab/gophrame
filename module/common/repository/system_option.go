package repository

import (
	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/common/domain"

	"gorm.io/gorm"
)

type SysOptionRepository struct {
	*gorm.DB `inject:"database"`
}

var sysOptionRepository = &SysOptionRepository{}

func init() {
	inject.InjectValue("sysOptionRepository", sysOptionRepository)
}

func (r *SysOptionRepository) GetDefaultOptions() (*domain.SysOptions, error) {
	result := &domain.SysOptions{TenantId: "DEFAULT", Options: make(map[string]*domain.SysOption)}

	var sysOptions = make([]*domain.SysOption, 0)
	if res := r.Where("tenant_id=?", "DEFAULT").Find(&sysOptions); res.Error == nil && res.RowsAffected > 0 {
		for _, option := range sysOptions {
			result.Options[option.Name] = option
		}
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SysOptionRepository) GetSystemOptions() (*domain.SysOptions, error) {
	result := &domain.SysOptions{TenantId: "SYSTEM", Options: make(map[string]*domain.SysOption)}

	var sysOptions = make([]*domain.SysOption, 0)
	if res := r.Where("tenant_id=?", "SYSTEM").Find(&sysOptions); res.Error == nil && res.RowsAffected > 0 {
		for _, option := range sysOptions {
			result.Options[option.Name] = option
		}
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SysOptionRepository) GetTenantOptions(tenantId string) (*domain.SysOptions, error) {
	result := &domain.SysOptions{TenantId: tenantId, Options: make(map[string]*domain.SysOption)}

	var sysOptions = make([]*domain.SysOption, 0)
	if res := r.Where("tenant_id=?", tenantId).Find(&sysOptions); res.Error == nil && res.RowsAffected > 0 {
		for _, option := range sysOptions {
			result.Options[option.Name] = option
		}
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SysOptionRepository) RemoveAllTenantOptions(tenantId string) error {
	return r.Delete(&domain.SysOption{}, "tenant_id=?", tenantId).Error
}
