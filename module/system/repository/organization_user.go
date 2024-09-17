package repository

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/system/domain"

	"gorm.io/gorm"
)

type OrganizationUserRepository struct {
	*gorm.DB `inject:"database"`
}

var organizationUserRepository *OrganizationUserRepository = &OrganizationUserRepository{}

func init() {
	inject.InjectValue("organizationUserRepository", organizationUserRepository)
}

// 查询类
func (a *OrganizationUserRepository) GetCount(organizationId int64, userName string) (count int64) {
	sql := `
		SELECT 
			count(*) as counts  
		FROM  
			sys_organization_user a, sys_user  b   
		WHERE 
			a.user_id=b.id 
			AND (a.organization_id=? or 0=?)
			AND (b.login LIKE ? or b.name like ?)
	`
	a.Raw(sql, organizationId, organizationId, "%"+userName+"%", "%"+userName+"%").First(&count)
	return
}

func (a *OrganizationUserRepository) ListMembers(organizationId int64, userName string, pageable query.Pageable) (count int64, data []domain.OrganizationMember) {
	count = a.GetCount(organizationId, userName)
	sql := `
		SELECT  
			c.id AS organization_id, 
			c.name AS organization_name, 
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
			sys_organization_user a
		JOIN sys_user b ON b.id = a.user_id
		JOIN sys_organization c ON c.id = a.organization_id
		LEFT JOIN sys_position d ON a.position_id=d.id
		WHERE  
			(a.organization_id=? or 0=?)
			AND (b.login LIKE ? or b.name LIKE ?)
		ORDER BY CONVERT(b.name USING GBK)
		LIMIT ?,?
	`
	a.Raw(sql, organizationId, organizationId, "%"+userName+"%", "%"+userName+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

func (a *OrganizationUserRepository) List(organizationId, userName string, pageable query.Pageable) (count int64, data []domain.OrganizationUser) {
	sql := `
		SELECT  
			a.*
		FROM  
			sys_organization_user a, sys_user b, sys_organization c
		WHERE  
			a.user_id=b.id AND c.id=a.organization_id
			AND (a.organization_id=? or 0=?)
			AND b.name LIKE ?
		LIMIT ?,?
	`
	a.Raw(sql, organizationId, organizationId, "%"+userName+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&data)
	return
}

// 新增
func (a *OrganizationUserRepository) InsertData(data *domain.OrganizationUser) bool {
	var counts int64
	if res := a.Model(&domain.OrganizationUser{}).Where("organization_id=? AND user_id=?", data.OrganizationId, data.UserId).Count(&counts); res.Error == nil && counts == 0 {
		if res := a.Create(data); res.Error == nil {
			return true
		} else {
			logger.Error("OrganizationUserRepository 新增失败", res.Error.Error())
		}
	} else {
		logger.Warn("OrganizationUserRepository 不允许重复新增")
	}
	return false
}

// 修改
func (a *OrganizationUserRepository) UpdateData(data *domain.OrganizationUser) bool {
	// Omit 表示忽略指定字段(CreatedTime)，其他字段全量更新
	if res := a.Omit("CreatedTime").Save(data); res.Error == nil {
		return true
	} else {
		logger.Error("OrganizationUserRepository 数据更新出错：", res.Error.Error())
	}
	return false
}

// 删除
func (a *OrganizationUserRepository) DeleteData(organizationId float64, userId string) bool {
	// 只能删除除了 admin 之外的用户
	var count int64
	a.Model(&domain.OrganizationUser{}).Select("user_id").Where("organization_id=? AND user_id=?", organizationId, userId).First(&count)
	if count < 1 {
		return true
	}

	if res := a.Where("organization_id=? AND user_id=?", organizationId, userId).Delete(&domain.OrganizationUser{}); res.Error == nil {
		return true
	} else {
		logger.Error("OrganizationUserRepository 删除数据出错：", res.Error.Error())
	}
	return false
}

// 修改
func (a *OrganizationUserRepository) GetByUserId(user_id string) (result []domain.OrganizationUser) {
	a.Where("user_id = ?", user_id).Find(&result)
	return
}
