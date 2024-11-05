package repository

import (
	"errors"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/operation/domain"

	"gorm.io/gorm"
)

type MenuRepository struct {
	*gorm.DB `inject:"database"`
}

var menuRepository *MenuRepository = &MenuRepository{}

func init() {
	inject.InjectValue("menuRepository", menuRepository)
}

func (a *MenuRepository) getCounts(fid int64, title string) (count int64) {
	var tx = a.Model(&domain.Menu{})
	if fid >= 0 {
		tx.Where("fid = ?", fid)
	}
	if title != "" {
		tx.Where("title LIKE ?", "%"+title+"%")
	}
	tx.Count(&count)
	return
}

// 查询
func (a *MenuRepository) List(fid int64, title string, pageable query.Pageable) (counts int64, data []*domain.Menu) {
	tx := a.Model(&domain.Menu{}).
		Order("sort ASC").
		Order("fid ASC").
		Limit(pageable.GetLimit()).
		Offset(pageable.GetOffset())

	if fid >= 0 {
		tx.Where("fid=?", fid)
	}

	if title != "" {
		tx.Where("title like ?", "%"+title+"%")
	}

	if pageable.NoCount() {
		if res := tx.Find(&data); res.Error == nil {
			return -1, data
		}
	} else {
		counts = a.getCounts(fid, title)
		if counts <= 0 {
			return 0, []*domain.Menu{}
		}
		if res := tx.Find(&data); res.Error == nil {
			return counts, data
		}
	}
	return 0, []*domain.Menu{}
}

// 查询
func (a *MenuRepository) ListWithButtons(fid int64, title string, pageable query.Pageable) (counts int64, data []*domain.Menu) {
	tx := a.Model(&domain.Menu{}).Preload("Buttons").
		Order("sort ASC").
		Order("fid ASC").
		Limit(pageable.GetLimit()).
		Offset(pageable.GetOffset())

	if fid >= 0 {
		tx.Where("fid=?", fid)
	}

	if title != "" {
		tx.Where("title like ?", "%"+title+"%")
	}

	if pageable.NoCount() {
		if res := tx.Find(&data); res.Error == nil {
			return -1, data
		}
	} else {
		counts = a.getCounts(fid, title)
		if counts <= 0 {
			return 0, []*domain.Menu{}
		}
		if res := tx.Find(&data); res.Error == nil {
			return counts, data
		}
	}
	return 0, []*domain.Menu{}
}

// 通过fid查询子节点数据
func (a *MenuRepository) GetById(id int64) (data *domain.Menu, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_menu WHERE fid=a.id) AS leaf
		FROM
			auth_menu a  
		WHERE
			id=?
	`
	err = a.Raw(sql, id).Scan(&data).Error
	return
}

// 通过fid查询子节点数据
func (a *MenuRepository) GetByFid(fid int64) (data []*domain.Menu, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_menu WHERE fid=a.id) AS leaf
		FROM
			auth_menu a  
		WHERE
			fid=?
	`
	err = a.Raw(sql, fid).Scan(&data).Error
	return
}

// 获取菜单fid的节点深度
func (a *MenuRepository) GetMenuLevel(fid int64) (level int64) {
	_ = a.Model(&domain.Menu{}).Select("level").Where("id=?", fid).First(&level)
	return
}

// 新增
func (a *MenuRepository) CreateMenu(data *domain.Menu) (bool, error) {
	var counts int64
	if res := a.Model(&domain.Menu{}).Where("fid=? AND (title=? OR name=?)", data.Fid, data.Title, data.Name).Count(&counts); res.Error == nil && counts == 0 {
		if res := a.Create(data); res.Error == nil {
			//新增菜单后,处理按钮
			go a.updatePathInfoNodeLevel(data.Id)
			return true, nil
		} else {
			logger.Error("MenuRepository 新增失败", res.Error.Error())
			return false, res.Error
		}
	} else {
		logger.Warn("MenuRepository 不允许重复新增")
		return false, errors.New("数据重复")
	}
}

// 更新
func (a *MenuRepository) UpdateMenu(data *domain.Menu) (bool, error) {
	// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
	if res := a.Save(data); res.Error == nil {
		go a.updatePathInfoNodeLevel(data.Id)
		return true, nil
	} else {
		logger.Error("MenuRepository 数据更新出错：", res.Error.Error())
		return false, res.Error
	}
}

// 新增、更新继续hook，更新path_info 、node_level 字段
func (a *MenuRepository) updatePathInfoNodeLevel(curItemid int64) bool {
	sql := `
		UPDATE 
			auth_menu a  
		LEFT JOIN auth_menu b ON a.fid=b.id
		SET 
			a.level=IFNULL(b.level,0)+1, 
			a.path_info=CONCAT(IFNULL(b.path_info,0),',',a.id)
		WHERE 
			a.id=?
		`
	if res := a.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		logger.Error("auth_menu 更新path_info失败", res.Error.Error())
	}
	return false
}

// 根据id查询是否有子节点数据
func (a *MenuRepository) GetSubNodeCount(id int64) (count int64) {
	if res := a.Model(&domain.Menu{}).Where("fid = ?", id).Count(&count); res.Error != nil {
		logger.Error("AuthSystemMenuModel 查询子节点是否有数据出错：", res.Error.Error())
	}
	return count
}

// 删除数据
func (a *MenuRepository) DeleteData(id int64) (bool, error) {
	if res := a.Delete(&domain.Menu{}, "id=?", id); res.Error == nil {
		go a.DeleteDataHook(id) // 删除菜单关联的所有数据
		return true, nil
	} else {
		logger.Error("MenuRepository 数据删除失败", res.Error.Error())
		return false, res.Error
	}
}

// 根据IDS获取菜单信息
func (a *MenuRepository) GetByIds(ids []int64) (result []*domain.Menu) {
	sql := `
			SELECT 
				a.id, a.fid, a.title, a.name, TRIM(a.icon) as icon, a.path, a.level, a.component, a.out_page,
				IFNULL((SELECT 0 FROM auth_menu b WHERE b.fid=a.id  LIMIT 1),1) as leaf
			FROM 
				auth_menu a 
			WHERE 
				id IN (?) 
				AND a.status=1 
			ORDER BY a.sort desc
		`
	a.Raw(sql, ids).Scan(&result)
	return
}

// 菜单主表数据删除，菜单关联的业务数据表同步删除
func (a *MenuRepository) DeleteDataHook(menuId int64) {
	//4. 删除菜单关联的待分配按钮子表
	sql := `DELETE FROM auth_button WHERE fid  = ? `
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("Repository 删除 auth_button 失败", res.Error.Error())
	}
}
