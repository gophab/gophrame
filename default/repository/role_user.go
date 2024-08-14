package repository

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/default/domain"

	"gorm.io/gorm"
)

type RoleUserRepository struct {
	*gorm.DB `inject:"database"`
}

var roleUserRepository *RoleUserRepository = &RoleUserRepository{}

func init() {
	inject.InjectValue("roleUserRepository", roleUserRepository)
}

// 查询类
func (a *RoleUserRepository) GetCount(roleId, search, tenantId string) (count int64) {
	sql := `
		SELECT 
			count(*) as counts  
		FROM  
			sys_role_user a, sys_user  b   
		WHERE 
			a.user_id=b.id 
			AND a.role_id=?
			AND (b.login LIKE ? or b.name like ? or b.email like ? or b.mobile like ?)
	`
	if tenantId != "" {
		sql += ` AND b.tenant_id = ?`
	} else {
		sql += ` AND b.tenant_id <> ?`
	}
	a.Raw(sql, roleId, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", tenantId).First(&count)
	return
}

func (a *RoleUserRepository) ListMembers(roleId, search, tenantId string, pageable query.Pageable) (count int64, data []*domain.RoleMember) {
	count = a.GetCount(roleId, search, tenantId)
	sql := `
		SELECT  
			c.id AS role_id, 
			c.name AS role_name, 
			b.id AS user_id,
			b.login,
			b.mobile,
			b.email,
			b.name,
			a.status,
			a.remark,
			d.title,
			d.name AS posistion_name,
			a.created_time,
			a.last_modified_time  
		FROM    
			sys_role_user a
		JOIN sys_user b ON b.id = a.user_id
		JOIN sys_role c ON c.id = a.role_id
		LEFT JOIN sys_position d ON a.position_id=d.id
		WHERE  
			a.role_id=?
			AND (b.login LIKE ? or b.name LIKE ? or b.email LIKE ? or b.mobile LIKE ?)
		ORDER BY CONVERT(b.name USING GBK)
		LIMIT ?,?
	`
	a.Raw(sql, roleId, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

func (a *RoleUserRepository) List(roleId, search, tenantId string, pageable query.Pageable) (count int64, data []*domain.RoleUser) {
	count = a.GetCount(roleId, search, tenantId)

	sql := `
		SELECT  
			a.*
		FROM  
			sys_role_user a, sys_user b, sys_role c
		WHERE  
			a.user_id=b.id AND c.id=a.role_id
			AND a.role_id=?
			AND (b.login LIKE ? or b.name LIKE ? or b.email LIKE ? or b.mobile LIKE ?)
	`

	if tenantId != "" {
		sql += ` AND b.tenant_id = ?`
	} else {
		sql += ` AND b.tenant_id <> ?`
	}

	sql += `
		ORDER BY CONVERT(b.name using GBK)
		LIMIT ?,?
	`
	a.Raw(sql, roleId, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", tenantId, pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

func (a *RoleUserRepository) ListUsers(roleId, search, tenantId string, pageable query.Pageable) (count int64, data []*domain.User) {
	count = a.GetCount(roleId, search, tenantId)

	sql := `
		SELECT  
			b.*
		FROM  
			sys_role_user a, sys_user b, sys_role c
		WHERE  
			a.user_id=b.id AND c.id=a.role_id
			AND a.role_id=?
			AND (b.login LIKE ? or b.name LIKE ? or b.email LIKE ? or b.mobile LIKE ?)
	`

	if tenantId != "" {
		sql += ` AND b.tenant_id = ?`
	} else {
		sql += ` AND b.tenant_id <> ?`
	}

	sql += `
		ORDER BY CONVERT(b.name using GBK)
		LIMIT ?,?
	`
	a.Raw(sql, roleId, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", tenantId, pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

func (a *RoleUserRepository) ListRoles(userId string, pageable query.Pageable) (count int64, data []*domain.Role) {
	sql := `
		SELECT  
			count(*)
		FROM  
			sys_role_user a
		WHERE  
			a.user_id=?
	`
	a.Raw(sql, userId).First(&count)

	sql = `
		SELECT  
			c.*
		FROM  
			sys_role_user a, sys_user b, sys_role c
		WHERE  
			a.user_id=b.id AND c.id=a.role_id
			AND a.user_id=?
		LIMIT ?,?
	`
	a.Raw(sql, userId, pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

// 新增
func (a *RoleUserRepository) InsertData(data *domain.RoleUser) bool {
	var counts int64
	if res := a.Model(&domain.RoleUser{}).Where("role_id=? AND user_id=?", data.RoleId, data.UserId).Count(&counts); res.Error == nil && counts == 0 {
		if res := a.Create(data); res.Error == nil {
			return true
		} else {
			logger.Error("RoleUserRepository 新增失败", res.Error.Error())
		}
	} else {
		logger.Warn("RoleUserRepository 不允许重复新增")
	}
	return false
}

// 修改
func (a *RoleUserRepository) UpdateData(data *domain.RoleUser) bool {
	// Omit 表示忽略指定字段(CreatedTime)，其他字段全量更新
	if res := a.Save(data); res.Error == nil {
		return true
	} else {
		logger.Error("RoleUserRepository 数据更新出错：", res.Error.Error())
	}
	return false
}

// 删除
func (a *RoleUserRepository) DeleteData(roleId string, userId string) bool {
	// 只能删除除了 admin 之外的用户
	res := a.Model(&domain.RoleUser{}).Select("user_id").Where("role_id=? AND user_id=?", roleId, userId).First(&userId)
	if res.RowsAffected < 1 {
		return true
	}

	if res := a.Where("role_id=? AND user_id=?", roleId, userId).Delete(&domain.RoleUser{}); res.Error == nil {
		return true
	} else {
		logger.Error("RoleUserRepository 删除数据出错：", res.Error.Error())
	}
	return false
}

// 修改
func (a *RoleUserRepository) GetByUserId(userId string) (result []*domain.RoleUser) {
	a.Where("user_id = ?", userId).Find(&result)
	return
}

// 修改
func (a *RoleUserRepository) GetByRoleId(roleId string) (result []*domain.RoleUser) {
	a.Where("role_id = ?", roleId).Find(&result)
	return
}
