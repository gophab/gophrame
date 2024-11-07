package repository

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/security/server"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/system/domain"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	server.IUserHandler
}

type RoleRepository struct {
	*gorm.DB `inject:"database"`
}

var roleRepository *RoleRepository = &RoleRepository{}

func init() {
	inject.InjectValue("roleRepository", roleRepository)
}

func (r *RoleRepository) ExistById(id string) (bool, error) {
	var role domain.Role
	err := r.Select("id").Where("id = ? AND del_flag = false ", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if role.Id != "" {
		return true, nil
	}

	return false, nil
}

func (r *RoleRepository) GetById(id string) (*domain.Role, error) {
	var role domain.Role
	err := r.Model(&domain.Role{}).Where("id = ? AND del_flag = false ", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetByIds(ids []string) (result []domain.Role) {
	r.Where("id IN ?", ids).Find(&result)
	return
}

func (r *RoleRepository) GetRoleTotal(maps interface{}) (int64, error) {
	var count int64
	if err := r.Model(&domain.Role{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RoleRepository) GetRolesAll() ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.Model(&domain.Role{}).Find(&roles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return roles, nil
}

func (r *RoleRepository) FindRoles(maps map[string]interface{}, pageable query.Pageable) (int64, []*domain.Role, error) {
	var role []*domain.Role
	tx := r.Model(&domain.Role{}).Where(maps)

	var count int64 = -1
	if !pageable.NoCount() {
		_ = tx.Count(&count)
	}

	if count != 0 {
		if err := query.Page(tx, pageable).Find(&role).Error; err != nil {
			return 0, nil, err
		}
	}

	return count, role, nil
}

func (r *RoleRepository) FindAvailableRoles(maps map[string]interface{}, pageable query.Pageable) (int64, []*domain.Role, error) {
	var role []*domain.Role

	tenantId := maps["tenant_id"]
	if tenantId != nil && tenantId.(string) != "" {
		delete(maps, "tenant_id")
	}

	tx := r.Model(&domain.Role{}).Where(maps)
	if tenantId != nil && tenantId.(string) != "" {
		tx.Where("tenant_id = ? or (tenant_id = 'SYSTEM' and scope = 'PUBLIC')", tenantId)
	}

	var count int64 = -1
	if !pageable.NoCount() {
		_ = tx.Count(&count)
	}

	if count != 0 {
		if err := query.Page(tx, pageable).Find(&role).Error; err != nil {
			return 0, nil, err
		}
	}

	return count, role, nil
}

func (r *RoleRepository) GetRoles(maps map[string]interface{}) ([]*domain.Role, error) {
	var roles []*domain.Role
	tx := r.Model(&domain.Role{})

	for k, v := range maps {
		if k == "user_id" {
			query := r.Model(&domain.RoleUser{}).Select("role_id").Where("user_id", v)
			tx.Where("id in (?)", query)
		} else {
			tx.Where(k+"=?", v)
		}
	}

	err := tx.Find(&roles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return roles, nil
}

func (r *RoleRepository) CheckRoleName(name string, tenantId string) (bool, error) {
	var role domain.Role
	err := r.Where("name = ? AND tenant_id = ? and del_flag = false ", name, tenantId).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if role.Id != "" {
		return true, nil
	}

	return false, nil
}

func (r *RoleRepository) CheckRoleNameId(name string, id string, tenantId string) (bool, error) {
	var role domain.Role
	err := r.Where("name = ? AND tenant_id = ? AND id != ? AND del_flag = false ", name, tenantId, id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if role.Id != "" {
		return true, nil
	}

	return false, nil
}

func (r *RoleRepository) PatchRole(id string, data map[string]interface{}) (*domain.Role, error) {
	data["id"] = id
	if err := r.Model(&domain.Role{}).Where("id = ? AND del_flag = false", id).UpdateColumns(util.DbFields(data)).Error; err != nil {
		return nil, err
	}

	return r.GetById(id)
}

func (r *RoleRepository) CreateRole(role *domain.Role) (*domain.Role, error) {
	if err := r.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) UpdateRole(role *domain.Role) (*domain.Role, error) {
	if err := r.Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) DeleteById(id string) error {
	if err := r.Model(&domain.Role{}).Where("id = ?", id).UpdateColumn("del_flag", true).Error; err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) CleanAllRole() error {
	if err := r.Unscoped().Where("del_flag = false ").Delete(&domain.Role{}).Error; err != nil {
		return err
	}

	return nil
}
