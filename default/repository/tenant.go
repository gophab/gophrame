package repository

import (
	"errors"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/default/domain"

	"gorm.io/gorm"
)

var tenantRepository *TenantRepository = &TenantRepository{}

func init() {
	inject.InjectValue("tenantRepository", tenantRepository)
}

type TenantRepository struct {
	*gorm.DB `inject:"database"`
}

func (r *TenantRepository) GetById(id string) (*domain.Tenant, error) {
	var result domain.Tenant
	if err := r.Model(&domain.Tenant{}).Where("id = ?", id).Find(&result); err.Error != nil {
		return nil, err.Error
	} else if err.RowsAffected == 0 {
		return nil, nil
	} else {
		return &result, nil
	}
}

func (a *TenantRepository) GetByIds(ids []string) (result []*domain.Tenant, err error) {
	err = a.Where("id IN ?", ids).Find(&result).Error
	return
}

// 新增
func (r *TenantRepository) CreateTenant(tenant *domain.Tenant) (*domain.Tenant, error) {
	var counts int64

	// 排除重名或者营业执照号重
	var tx = r.Model(&domain.Tenant{})
	if tenant.Name != "" && tenant.LicenseId != "" {
		tx.Where("name=? or license_id=?", tenant.Name, tenant.LicenseId)
	} else {
		tx.Where("name=?", tenant.Name)
	}

	if res := tx.Count(&counts); res.Error == nil && counts > 0 {
		return nil, errors.New("tenant 重复")
	}

	if res := r.Create(tenant); res.Error == nil {
		return tenant, nil
	} else {
		logger.Error("Tenant 数据新增出错：", res.Error.Error())
		return nil, res.Error
	}
}

// 更新
func (r *TenantRepository) UpdateTenant(tenant *domain.Tenant) (bool, error) {
	var counts int64

	// 同一个地区下不存在相同名称的区域
	var tx = r.Model(&domain.Tenant{})
	if tenant.Name != "" || tenant.LicenseId != "" {
		if tenant.Name != "" && tenant.LicenseId != "" {
			tx.Where("id <> ? and (name=? or license_id=?)", tenant.Id, tenant.Name, tenant.LicenseId)
		} else if tenant.Name != "" {
			tx.Where("id <> ? and name=?", tenant.Id, tenant.Name)
		} else if tenant.LicenseId != "" {
			tx.Where("id <> ? and license_id=?", tenant.Id, tenant.LicenseId)
		}

		if res := tx.Count(&counts); res.Error == nil && counts > 0 {
			return false, errors.New("tenant 重复")
		}
	}

	// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
	if res := r.Omit("CreatedTime").Save(tenant); res.Error == nil {
		return true, nil
	} else {
		logger.Error("Tenant 数据更新失败，错误详情：", res.Error.Error())
		return false, res.Error
	}
}

// 删除
func (r *TenantRepository) DeleteById(id string) bool {
	err := r.Model(&domain.Tenant{}).Where("id=?", id).Update("del_flag", true).Error
	if err == nil {
		return true
	} else {
		logger.Error("Tenant 删除数据出错：", err.Error())
	}
	return false
}

func (r *TenantRepository) Find(conds map[string]interface{}, pageable query.Pageable) (total int64, list []*domain.Tenant) {
	var tx = r.DB.Model(&domain.Tenant{})

	var search = conds["search"]
	var id = conds["id"]
	var name = conds["name"]
	var licenseId = conds["license_id"]

	if search != nil && search != "" {
		tx = tx.Where("name like ? or name_cn like ? or name_tw like ? or name_en like ? or license_id like ? or id = ?",
			"%"+search.(string)+"%",
			"%"+search.(string)+"%",
			"%"+search.(string)+"%",
			"%"+search.(string)+"%",
			"%"+search.(string)+"%",
			search)
	} else {
		if name != nil && name != "" {
			tx = tx.Where("name like ?", "%"+name.(string)+"%")
		}
		if licenseId != nil && licenseId != "" {
			tx = tx.Where("license_id like ?", "%"+licenseId.(string)+"%")
		}
		if id != nil && id != "" {
			tx = tx.Where("id = ?", id)
		}
	}

	list = make([]*domain.Tenant, 0)
	total = 0

	if tx.Count(&total).Error != nil || total == 0 {
		return
	}

	query.Page(tx, pageable).Find(&list)
	return
}
