package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
)

var organizationService *OrganizationService = &OrganizationService{}

func GetOrganizationService() *OrganizationService {
	return organizationService
}

func init() {
	inject.InjectValue("organizationService", organizationService)
}

type OrganizationService struct {
	service.BaseService
	OrganizationRepository *repository.OrganizationRepository `inject:"organizationRepository"`
}

func (s *OrganizationService) CreateOrganization(organization *domain.Organization) (*domain.Organization, error) {
	if b, err := s.OrganizationRepository.InsertData(organization); b && err == nil {
		return organization, nil
	} else {
		return nil, err
	}
}

func (s *OrganizationService) UpdateOrganization(organization *domain.Organization) (*domain.Organization, error) {
	if b, err := s.OrganizationRepository.UpdateData(organization); b && err == nil {
		return organization, nil
	} else {
		return nil, err
	}
}

func (s *OrganizationService) DeleteOrganization(id string) (bool, error) {
	result := s.OrganizationRepository.DeleteData(id)
	return result, nil
}

func (s *OrganizationService) GetById(id string) (*domain.Organization, error) {
	return s.OrganizationRepository.GetById(id)
}

func (s *OrganizationService) GetByIds(ids []string) ([]*domain.Organization, error) {
	return s.OrganizationRepository.GetByIds(ids)
}

func (s *OrganizationService) List(fid string, name string, pageable query.Pageable) (total int64, list []*domain.Organization) {
	total, list = s.OrganizationRepository.List(fid, name, pageable)

	return total, list
}

func (s *OrganizationService) GetSubList(fid string) []*domain.Organization {
	return s.OrganizationRepository.GetSubListByFid(fid)
}

func (s *OrganizationService) GetParentList(id string) []*domain.Organization {
	list, _ := s.OrganizationRepository.GetAllParentsById(id)
	return list
}

func (s *OrganizationService) HasSubNode(fid string) int64 {
	return s.OrganizationRepository.HasSubNode(fid)
}
