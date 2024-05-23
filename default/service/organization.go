package service

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/service"

	"github.com/wjshen/gophrame/default/domain"
	"github.com/wjshen/gophrame/default/repository"
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

func (s *OrganizationService) DeleteOrganization(id int64) (bool, error) {
	result := s.OrganizationRepository.DeleteData(id)
	return result, nil
}

func (s *OrganizationService) GetById(id int64) (*domain.Organization, error) {
	return s.OrganizationRepository.GetById(id)
}

func (s *OrganizationService) List(fid int64, name string, pageable query.Pageable) (total int64, list []domain.Organization) {
	total, list = s.OrganizationRepository.List(fid, name, pageable)

	return total, list
}

func (s *OrganizationService) GetSubList(fid int64) []domain.Organization {
	return s.OrganizationRepository.GetSubListByfid(fid)
}

func (s *OrganizationService) HasSubNode(fid int64) int64 {
	return s.OrganizationRepository.HasSubNode(fid)
}
