package service

import (
	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/service"
)

type ContentTemplateService struct {
	service.BaseService
	ContentTemplateRepository *repository.ContentTemplateRepository `inject:"contentTemplateRepository"`
}

var contentTemplateService = &ContentTemplateService{}

func init() {
	inject.InjectValue("contentTemplateService", contentTemplateService)
}

func (s *ContentTemplateService) GetById(id string) (*domain.ContentTemplate, error) {
	return s.ContentTemplateRepository.GetById(id)
}

func (s *ContentTemplateService) GetByTypeAndSceneAndTenantId(typeName, scene, tenantId string) (*domain.ContentTemplate, error) {
	return s.ContentTemplateRepository.GetByTypeAndSceneAndTenantId(typeName, scene, tenantId)
}

func (s *ContentTemplateService) FindAll(conds map[string]any, pageable query.Pageable) (int64, []*domain.ContentTemplate, error) {
	return s.ContentTemplateRepository.FindAll(conds, pageable)
}

func (s *ContentTemplateService) CreateContentTemplate(template *domain.ContentTemplate) (*domain.ContentTemplate, error) {
	return s.ContentTemplateRepository.CreateContentTemplate(template)
}

func (s *ContentTemplateService) UpdateContentTemplate(template *domain.ContentTemplate) (*domain.ContentTemplate, error) {
	return s.ContentTemplateRepository.CreateContentTemplate(template)
}

func (s *ContentTemplateService) PatchContentTemplate(id string, data map[string]any) (result *domain.ContentTemplate, err error) {
	return s.ContentTemplateRepository.PatchContentTemplate(id, data)
}

func (s *ContentTemplateService) DeleteContentTemplate(id string) error {
	return s.ContentTemplateRepository.DeleteById(id)
}

func (s *ContentTemplateService) GetContentTemplate(typeName, scene string) (title, content string) {
	contentTemplate, err := s.ContentTemplateRepository.GetByTypeAndSceneAndTenantId(typeName, scene, SecurityUtil.GetCurrentTenantId(nil))
	if err == nil && contentTemplate != nil {
		return contentTemplate.Title, contentTemplate.Content
	}
	return scene, scene
}
