package repository

import (
	"errors"
	"fmt"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/operation/domain"

	"gorm.io/gorm"
)

type ModuleRepository struct {
	*gorm.DB `inject:"database"`
}

var moduleRepository *ModuleRepository = &ModuleRepository{}

func init() {
	inject.InjectValue("moduleRepository", moduleRepository)
}

// 查询
func (a *ModuleRepository) List(fid int64, title string, tenantId string, pageable query.Pageable) (counts int64, data []*domain.Module) {
	tx := a.Model(&domain.Module{}).
		Order("fid ASC").
		Limit(pageable.GetLimit()).
		Offset(pageable.GetOffset())

	if fid >= 0 {
		tx.Where("fid=?", fid)
	}

	if title != "" {
		tx.Where("title like ?", "%"+title+"%")
	}

	if tenantId != "" {
		tx.Where("tenant_id = ?", tenantId)
	}

	if pageable.NoCount() {
		if res := tx.Find(&data); res.Error == nil {
			return -1, data
		}
	} else {
		if res := tx.Count(&counts); res.Error != nil || counts <= 0 {
			return 0, []*domain.Module{}
		}
		if res := tx.Find(&data); res.Error == nil {
			return counts, data
		}
	}
	return 0, []*domain.Module{}
}

// 通过fid查询子节点数据
func (a *ModuleRepository) GetById(id int64) (data *domain.Module, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_module WHERE fid=a.id) AS leaf
		FROM
			auth_module a  
		WHERE
			id=?
	`
	err = a.Raw(sql, id).Scan(&data).Error
	return
}

// 通过fid查询子节点数据
func (a *ModuleRepository) GetByFid(fid int64) (data []*domain.Module, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_module WHERE fid=a.id) AS leaf
		FROM
			auth_module a  
		WHERE
			fid=?
	`
	err = a.Raw(sql, fid).Scan(&data).Error
	return
}

// 新增
func (a *ModuleRepository) CreateModule(data *domain.Module) (bool, error) {
	var counts int64
	if res := a.Model(&domain.Module{}).Where("fid=? AND (title=? OR name=?)", data.Fid, data.Title, data.Name).Count(&counts); res.Error == nil && counts == 0 {
		if res := a.Create(data); res.Error == nil {
			//新增菜单后,处理按钮
			go a.updatePathInfoNodeLevel(data.Id)
			return true, nil
		} else {
			logger.Error("ModuleRepository 新增失败", res.Error.Error())
			return false, res.Error
		}
	} else {
		logger.Warn("ModuleRepository 不允许重复新增")
		return false, errors.New("数据重复")
	}
}

// 更新
func (a *ModuleRepository) UpdateModule(data *domain.Module) (bool, error) {
	// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
	if res := a.Save(data); res.Error == nil {
		go a.updatePathInfoNodeLevel(data.Id)
		return true, nil
	} else {
		logger.Error("ModuleRepository 数据更新出错：", res.Error.Error())
		return false, res.Error
	}
}

// 新增、更新继续hook，更新path_info 、node_level 字段
func (a *ModuleRepository) updatePathInfoNodeLevel(curItemid int64) bool {
	sql := `
		UPDATE 
			auth_module a  
		LEFT JOIN auth_module b ON a.fid=b.id
		SET
			a.fids=CONCAT(IFNULL(b.fids, 0), ',', a.id),
			a.path=CONCAT(CASE WHEN ISNULL(b.path) THEN '' ELSE b.path||'.' END, a.name)
		WHERE 
			a.id=?
		`
	if res := a.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		logger.Error("auth_module 更新path_info失败", res.Error.Error())
	}
	return false
}

// 删除数据
func (a *ModuleRepository) DeleteData(id int64) (bool, error) {
	if res := a.Delete(&domain.Module{}, id); res.Error == nil {
		// 删除下
		a.Delete(&domain.Module{}, "fids like ?", fmt.Sprintf("%%,%d,%%", id))
		return true, nil
	} else {
		logger.Error("ModuleRepository 数据删除失败", res.Error.Error())
		return false, res.Error
	}
}

// 根据IDS获取菜单信息
func (a *ModuleRepository) GetByIds(ids []int64) (result []*domain.Module, err error) {
	sql := `
			SELECT 
				a.*
			FROM 
				auth_module a 
			WHERE 
				id IN (?) 
		`
	err = a.Raw(sql, ids).Scan(&result).Error
	return
}

func (a *ModuleRepository) GetTenantModules(tenantId string) ([]*domain.Module, error) {
	var results = make([]*domain.Module, 0)
	if res := a.Model(&domain.Module{}).Where("tenant_id = ?", tenantId).Where("status = ?", 1).Find(&results); res.Error == nil {
		// var list = make([]*auth.Module, 0)
		// if err := a.makeTree(results, &list); err == nil {
		// 	return list, nil
		// }
		return results, nil
	} else {
		return nil, res.Error
	}
}

func (a *ModuleRepository) DeleteTenantModules(tenantId string) {
	if res := a.Delete(&domain.Module{}, "tenant_id = ?", tenantId); res.Error != nil {
		logger.Warn("Delete TenantModules error: ", res.Error.Error())
	}
}

// 给角色分配系统菜单、按钮
func (a *ModuleRepository) SetTenantModules(tenantId string, data []*domain.Module) {
	// 1. clear
	a.DeleteTenantModules(tenantId)

	// 2. add
	for _, module := range data {
		var m = &domain.Module{
			ModuleInfo: module.ModuleInfo,
		}
		m.Oid = module.Id
		m.TenantId = tenantId

		if res := a.Create(m); res.Error != nil {
			logger.Warn("Persist TenantModule error: ", res.Error.Error())
			break
		}
	}
}

func (a *ModuleRepository) GetTenantModule(module, tenantId string) (*domain.Module, error) {
	var result domain.Module
	if res := a.Model(&domain.Module{}).Where("path = ?", module).Where("tenant_id = ?", tenantId).First(&result); res.Error == nil && res.RowsAffected > 0 {

		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (a *ModuleRepository) makeTree(src []*domain.Module, dest *[]*domain.Module) error {
	var result = *dest
	var srcMap = make(map[int64]*domain.Module)
	for _, item := range src {
		srcMap[item.Id] = item
	}
	for _, item := range src {
		if item.Fid != 0 {
			var parent = srcMap[item.Fid]
			if parent != nil {
				if parent.Children == nil {
					parent.Children = make([]*domain.Module, 0)
				}
				parent.Children = append(parent.Children, item)
			}
		} else {
			result = append(result, item)
		}
	}
	*dest = result
	return nil
}
