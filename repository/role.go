package repository

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/security/server"

	"github.com/wjshen/gophrame/domain"

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

func (r *RoleRepository) ExistRoleByID(id string) (bool, error) {
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

func (r *RoleRepository) GetRoles(maps interface{}, pageable query.Pageable) ([]*domain.Role, error) {
	var role []*domain.Role
	err := query.Page(r.Preload("Menus").Where(maps), pageable).Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) GetRole(id string) (*domain.Role, error) {
	var role domain.Role
	err := r.Preload("Menus").Where("id = ? AND del_flag = false ", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &role, nil
}
func (r *RoleRepository) CheckRoleName(name string) (bool, error) {
	var role domain.Role
	err := r.Where("name = ? AND del_flag = false ", name).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if role.Id != "" {
		return true, nil
	}

	return false, nil
}

func (r *RoleRepository) CheckRoleNameId(name string, id string) (bool, error) {
	var role domain.Role
	err := r.Where("name = ? AND id != ? AND del_flag = false ", name, id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if role.Id != "" {
		return true, nil
	}

	return false, nil
}

func (r *RoleRepository) EditRole(id string, data map[string]interface{}) error {
	var role []domain.Role

	if err := r.Where("id = ? AND del_flag = false ", id).Find(&role).Error; err != nil {
		return err
	}
	r.Model(&role).UpdateColumns(data)

	return nil
}

func (r *RoleRepository) AddRole(data map[string]interface{}) (*domain.Role, error) {
	role := domain.Role{
		Name: data["name"].(string),
	}
	if err := r.Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) DeleteRole(id string) error {
	var role domain.Role
	r.Where("id = ?", id).Find(&role)
	if err := r.Where("id = ?", id).Delete(&role).Error; err != nil {
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

func (r *RoleRepository) GetRolesAll() ([]*domain.Role, error) {
	var role []*domain.Role
	err := r.Preload("Menus").Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return role, nil
}
