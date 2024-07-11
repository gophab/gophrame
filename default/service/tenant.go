package service

import (
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
)

var tenantService *TenantService = &TenantService{}

func GetTenantService() *TenantService {
	return tenantService
}

func init() {
	inject.InjectValue("tenantService", tenantService)
}

type TenantService struct {
	service.BaseService
	TenantRepository *repository.TenantRepository `inject:"tenantRepository"`
}

func (s *TenantService) Create(tenant *domain.Tenant) (*domain.Tenant, error) {
	if result, err := s.TenantRepository.CreateTenant(tenant); result != nil && err == nil {
		eventbus.PublishEvent("TENANT_CREATED", tenant)
		return result, nil
	} else {
		return nil, err
	}
}

func (s *TenantService) Update(tenant *domain.Tenant) (*domain.Tenant, error) {
	if b, err := s.TenantRepository.UpdateTenant(tenant); b && err == nil {
		eventbus.PublishEvent("TENANT_UPDATED", tenant)
		return tenant, nil
	} else {
		return nil, err
	}
}

func (s *TenantService) Patch(id string, column string, value interface{}) (*domain.Tenant, error) {
	if res := s.TenantRepository.Model(&domain.Tenant{}).Where("id=?", id).UpdateColumn(column, value); res.Error != nil {
		return nil, res.Error
	} else {
		if tenant, err := s.GetById(id); err == nil {
			eventbus.PublishEvent("TENANT_UPDATED", tenant)
			return tenant, err
		} else {
			return nil, err
		}
	}
}

func (s *TenantService) PatchAll(id string, kv map[string]interface{}) (*domain.Tenant, error) {
	if res := s.TenantRepository.Model(&domain.Tenant{}).Where("id=?", id).UpdateColumns(kv); res.Error != nil {
		return nil, res.Error
	}

	if tenant, err := s.GetById(id); err != nil {
		return nil, err
	} else {
		eventbus.PublishEvent("TENANT_UPDATED", tenant)
		return tenant, err
	}
}
func (s *TenantService) DeleteTenant(tenant *domain.Tenant) (bool, error) {
	result := s.TenantRepository.DeleteById(tenant.Id)
	eventbus.PublishEvent("TENANT_DELETED", tenant)
	return result, nil
}

func (s *TenantService) DeleteById(id string) (bool, error) {
	if tenant, err := s.GetById(id); err != nil {
		return s.DeleteTenant(tenant)
	} else {
		return true, nil
	}
}

func (s *TenantService) GetById(id string) (*domain.Tenant, error) {
	return s.TenantRepository.GetById(id)
}

func (s *TenantService) Find(name, licenseId string, pageable query.Pageable) (total int64, list []*domain.Tenant) {
	total, list = s.TenantRepository.Find(name, licenseId, pageable)
	return
}
