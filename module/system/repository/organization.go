package repository

import (
	"errors"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/system/domain"

	"gorm.io/gorm"
)

var organizationRepository *OrganizationRepository = &OrganizationRepository{}

func init() {
	inject.InjectValue("organizationRepository", organizationRepository)
}

type OrganizationRepository struct {
	*gorm.DB `inject:"database"`
}

func (r *OrganizationRepository) GetCount(fid string, name string) (count int64) {
	r.Model(&domain.Organization{}).Select("id").Where("fid=? AND name like ?", fid, "%"+name+"%").Count(&count)
	return
}

func (r *OrganizationRepository) GetById(id string) (*domain.Organization, error) {
	var result domain.Organization
	if err := r.Model(&domain.Organization{}).Where("id = ?", id).Find(&result); err.Error != nil {
		return nil, err.Error
	} else if err.RowsAffected == 0 {
		return nil, nil
	} else {
		return &result, nil
	}
}

func (r *OrganizationRepository) GetParentById(id string) (*domain.Organization, error) {
	if org, err := r.GetById(id); err == nil && org != nil {
		if org.Fid != "" {
			return r.GetById(org.Fid)
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

func (r *OrganizationRepository) GetAllParentsById(id string) ([]*domain.Organization, error) {
	if parent, err := r.GetParentById(id); err == nil && parent != nil {
		if parent.Fid != "" {
			if parents, err := r.GetAllParentsById(parent.Fid); err == nil && len(parents) > 0 {
				return append([]*domain.Organization{parent}, parents...), nil
			} else {
				return []*domain.Organization{parent}, err
			}
		} else {
			return []*domain.Organization{parent}, nil
		}
	} else {
		return []*domain.Organization{}, err
	}
}

// 查询
func (r *OrganizationRepository) List(fid string, name string, pageable query.Pageable) (counts int64, list []*domain.Organization) {
	if counts = r.GetCount(fid, name); counts > 0 {
		sql := `
			SELECT
				a.*
			FROM sys_organization a
			WHERE   a.fid= ? AND   a.name LIKE  ? ORDER  BY a.fid ASC, a.id  ASC
			LIMIT ? , ?
		`
		_ = r.Raw(sql, fid, "%"+name+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&list)
	}
	return
}

// 根据fid查询子级节点全部数据
func (r *OrganizationRepository) GetSubListByFid(fid string) []*domain.Organization {
	sql := `
		SELECT
			a.*,
			(SELECT  CASE  WHEN  COUNT(*)=0 THEN 1 ELSE 0 END FROM sys_organization WHERE fid=a.id) AS leaf
		FROM sys_organization a
		WHERE fid = ?
	`
	var inSlice []*domain.Organization
	if res := r.Raw(sql, fid).Find(&inSlice); res.Error == nil && len(inSlice) > 0 {
		return inSlice
	} else if res.Error != nil {
		logger.Error("Organization 根据fid查询子级出错:", res.Error.Error())
	}
	return nil
}

// 根据fid查询子级节点全部数据
func (r *OrganizationRepository) GetAllSubListByFid(fid string) []*domain.Organization {
	var subs = r.GetSubListByFid(fid)
	if len(subs) > 0 {
		for _, sub := range subs {
			list := r.GetAllSubListByFid(sub.Id)
			if len(list) > 0 {
				subs = append(subs, list...)
			}
		}
	}
	return subs
}

// 新增
func (r *OrganizationRepository) InsertData(organization *domain.Organization) (bool, error) {
	var counts int64

	// 同一个地区下不存在相同名称的区域
	if res := r.Model(&domain.Organization{}).Where("fid=? and name=?", organization.Fid, organization.Name).Count(&counts); res.Error == nil && counts > 0 {
		return false, errors.New("organization 重复")
	}

	if res := r.Create(*organization); res.Error == nil {
		_ = r.updatePathInfoNodeLevel(organization.Id)
		return true, nil
	} else {
		logger.Error("Organization 数据新增出错：", res.Error.Error())
		return false, res.Error
	}
}

// 更新
func (r *OrganizationRepository) UpdateData(organization *domain.Organization) (bool, error) {
	var counts int64

	// 同一个地区下不存在相同名称的区域
	if res := r.Model(&domain.Organization{}).Where("id <> ? and fid=? and name=?", organization.Id, organization.Fid, organization.Name).Count(&counts); res.Error == nil && counts > 0 {
		return false, errors.New("organization 重复")
	}

	// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
	if res := r.Omit("CreatedTime").Save(*organization); res.Error == nil {
		_ = r.updatePathInfoNodeLevel(organization.Id)
		return true, nil
	} else {
		logger.Error("Organization 数据更新失败，错误详情：", res.Error.Error())
		return false, res.Error
	}
}

// 删除
func (r *OrganizationRepository) DeleteData(id string) bool {
	if res := r.Delete(&domain.Organization{}, "id=?", id); res.Error == nil {
		return true
	} else {
		logger.Error("Organization 删除数据出错：", res.Error.Error())
	}
	return false
}

// 查询该 id 是否存在子节点
func (r *OrganizationRepository) HasSubNode(id string) (count int64) {
	r.Model(&domain.Organization{}).Select("id").Where("fid=?", id).Count(&count)
	return count
}

// 更新path_info 、node_level 字段
func (r *OrganizationRepository) updatePathInfoNodeLevel(curItemid string) bool {
	sql := `
		UPDATE sys_organization a  LEFT JOIN sys_organization  b
		ON  a.fid=b.id
		SET  a.node_level=b.node_level+1,  a.path_info=CONCAT(b.path_info,',',a.id)
		WHERE  a.id=?
		`
	if res := r.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		logger.Error("Organization 更新 node_level , path_info 失败", res.Error.Error())
	}
	return false
}

func (a *OrganizationRepository) GetByIds(ids []string) (result []*domain.Organization) {
	a.Where("id IN ?", ids).Find(&result)
	return
}
