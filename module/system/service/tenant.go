package service

import (
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
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

func (s *TenantService) Patch(id string, column string, value any) (*domain.Tenant, error) {
	if res := s.TenantRepository.Model(&domain.Tenant{}).Where("id=?", id).UpdateColumn(util.DbFieldName(column), value); res.Error != nil {
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

func (s *TenantService) PatchAll(id string, kv map[string]any) (*domain.Tenant, error) {
	kv["id"] = id

	if res := s.TenantRepository.Model(&domain.Tenant{}).Where("id=?", id).UpdateColumns(util.DbFields(kv)); res.Error != nil {
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
	if tenant, err := s.GetById(id); tenant != nil {
		return s.DeleteTenant(tenant)
	} else {
		return true, err
	}
}

func (s *TenantService) GetById(id string) (*domain.Tenant, error) {
	return s.TenantRepository.GetById(id)
}

func (s *TenantService) GetByIds(ids []string) ([]*domain.Tenant, error) {
	return s.TenantRepository.GetByIds(ids)
}

func (s *TenantService) Find(conds map[string]any, pageable query.Pageable) (total int64, list []*domain.Tenant) {
	total, list = s.TenantRepository.Find(util.DbFields(conds), pageable)
	return
}
