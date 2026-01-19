package repository

import (
	"github.com/gophab/gophrame/module/common/domain"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/transaction"
	"github.com/gophab/gophrame/core/util"
	"gorm.io/gorm"
)

type ContentTemplateRepository struct {
	*gorm.DB `inject:"database"`
}

var contentTemplateRepository = &ContentTemplateRepository{}

func init() {
	inject.InjectValue("contentTemplateRepository", contentTemplateRepository)
}

func (r *ContentTemplateRepository) GetById(id string) (*domain.ContentTemplate, error) {
	var result domain.ContentTemplate
	if res := transaction.Session().Model(&domain.ContentTemplate{}).
		Where("id=?", id).
		Find(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *ContentTemplateRepository) GetByTypeAndSceneAndTenantId(typeName, scene, tenantId string) (*domain.ContentTemplate, error) {
	var results = make([]*domain.ContentTemplate, 0)
	if res := transaction.Session().Model(&domain.ContentTemplate{}).
		Where("type=? and scene=? and (tenant_id=? or tenant_id='SYSTEM') and status<>0", typeName, scene, tenantId).
		Find(&results); res.Error == nil && res.RowsAffected > 0 {
		for _, result := range results {
			if result.TenantId != "SYSTEM" {
				return result, nil
			}
		}
		return results[0], nil
	} else {
		return nil, res.Error
	}
}

func (r *ContentTemplateRepository) GetAll(conds map[string]any) (result []*domain.ContentTemplate, err error) {
	q := transaction.Session().Model(&domain.ContentTemplate{})
	for k, v := range conds {
		q = q.Where(k+"=?", v)
	}

	if res := q.Find(&result); res.Error == nil && res.RowsAffected > 0 {
		return result, nil
	} else {
		return []*domain.ContentTemplate{}, res.Error
	}
}

func (r *ContentTemplateRepository) FindAll(conds map[string]any, pageable query.Pageable) (count int64, result []*domain.ContentTemplate, err error) {
	q := transaction.Session().Model(&domain.ContentTemplate{})
	for k, v := range conds {
		if k == "search" {
			q = q.Where("name like ?", "%"+v.(string)+"%")
		} else {
			q = q.Where(k+"=?", v)
		}
	}

	q.Count(&count)

	order := "name"
	switch r.Dialector.Name() {
	case "mysql":
		order = "convert(name using gbk)"
	case "postgres":
		order = "conver(name, 'utf-8', 'gbk')"
	}

	if res := query.Page(q, pageable).Order(order).Find(&result); res.Error == nil && res.RowsAffected > 0 {
		return count, result, nil
	} else {
		return 0, []*domain.ContentTemplate{}, res.Error
	}
}

func (r *ContentTemplateRepository) CreateContentTemplate(template *domain.ContentTemplate) (*domain.ContentTemplate, error) {
	if res := transaction.Session().Model(&domain.ContentTemplate{}).Create(template); res.Error == nil && res.RowsAffected > 0 {
		return template, nil
	} else {
		return nil, res.Error
	}
}

func (r *ContentTemplateRepository) UpdateContentTemplate(template *domain.ContentTemplate) (*domain.ContentTemplate, error) {
	if res := transaction.Session().Model(&domain.ContentTemplate{}).Save(template); res.Error == nil && res.RowsAffected > 0 {
		return template, nil
	} else {
		return nil, res.Error
	}
}

func (r *ContentTemplateRepository) PatchContentTemplate(id string, data map[string]any) (result *domain.ContentTemplate, err error) {
	if result, err = r.GetById(id); err != nil || result == nil {
		return
	}

	data["id"] = id
	if res := transaction.Session().Model(&domain.ContentTemplate{}).Where("id=?", id).UpdateColumns(util.DbFields(data)); res.Error == nil {
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (r *ContentTemplateRepository) DeleteById(id string) error {
	return transaction.Session().Delete(&domain.ContentTemplate{}, "id=?", id).Error
}
