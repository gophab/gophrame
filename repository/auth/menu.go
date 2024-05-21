package auth

import (
	"errors"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/domain/auth"

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
	sql := "SELECT COUNT(*) AS counts FROM `auth_menu` WHERE fid=? AND title LIKE ?"
	a.Raw(sql, fid, "%"+title+"%").First(&count)
	return
}

// 查询
func (a *MenuRepository) List(fid int64, title string, pageable query.Pageable) (counts int64, data []auth.Menu) {
	counts = a.getCounts(fid, title)
	if counts > 0 {
		sql := `
			SELECT 
				a.*
			FROM 
				auth_menu a  
			WHERE 
				fid=? 
				AND title LIKE ? 
			ORDER BY a.sort DESC, a.fid  ASC
			LIMIT ?,?
		`
		if res := a.Raw(sql, fid, "%"+title+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&data); res.Error == nil {
			return counts, data
		}
	}
	return 0, nil
}

// 查询
func (a *MenuRepository) ListWithButtons(fid int64, title string, pageable query.Pageable) (counts int64, data []auth.MenuWithButton) {
	counts = a.getCounts(fid, title)
	if counts > 0 {
		sql := `
			SELECT 
				a.*,
				IFNULL(b.menu_id,0) AS menu_id,
				IFNULL(c.id,0) AS button_id,
				IFNULL(c.cn_name,'') AS button_name, 
				IFNULL(c.color,'') AS button_color 
			FROM 
				auth_menu a  
			LEFT JOIN auth_menu_button b ON a.id =b.menu_id
			LEFT JOIN auth_button c ON c.id=b.button_id
			WHERE 
				a.id IN (
					SELECT id FROM (SELECT id FROM auth_menu WHERE fid=? AND title LIKE ? LIMIT ?,?) AS t_tmp 
				)
			ORDER BY a.sort DESC, a.fid  ASC, button_id ASC
		`
		result := []auth.MenuButton{}
		if res := a.Raw(sql, fid, "%"+title+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&result); res.Error == nil {
			if err := util.CreateSqlResFormatFactory().ScanToTreeData(result, &data); err == nil {
				return counts, data
			} else {
				logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
			}

		}
	}
	return 0, nil
}

// 通过fid查询子节点数据
func (a *MenuRepository) GetById(id int64) (data auth.Menu, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) >0 THEN 1 ELSE 0 END FROM auth_menu WHERE fid=a.id) AS has_sub_node,
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
func (a *MenuRepository) GetByFid(fid int64) (data []auth.Menu, err error) {
	sql := `
		SELECT  
			a.*,
			(SELECT CASE WHEN COUNT(*) >0 THEN 1 ELSE 0 END FROM auth_menu WHERE fid=a.id) AS has_sub_node,
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
func (a *MenuRepository) GetMenuLevel(fid int64) (nodeLevel int64) {
	_ = a.Model(&auth.Menu{}).Select("node_level").Where("id=?", fid).First(&nodeLevel)
	return
}

// 新增
func (a *MenuRepository) InsertData(data *auth.Menu) (bool, error) {
	var counts int64
	if res := a.Model(&auth.Menu{}).Where("fid=? AND (title=? OR name=?)", data.Fid, data.Title, data.Name).Count(&counts); res.Error == nil && counts == 0 {
		if res := a.Create(*data); res.Error == nil {
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
func (a *MenuRepository) UpdateData(data *auth.Menu) (bool, error) {
	// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
	if res := a.Omit("CreatedTime").Save(*data); res.Error == nil {
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
			a.node_level=IFNULL(b.node_level,0)+1, 
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
	if res := a.Model(&auth.Menu{}).Where("fid = ?", id).Count(&count); res.Error != nil {
		logger.Error("AuthSystemMenuModel 查询子节点是否有数据出错：", res.Error.Error())
	}
	return count
}

// 删除数据
func (a *MenuRepository) DeleteData(id int64) (bool, error) {
	if res := a.Delete(a, id); res.Error == nil {
		go a.DeleteDataHook(id) // 删除菜单关联的所有数据
		return true, nil
	} else {
		logger.Error("MenuRepository 数据删除失败", res.Error.Error())
		return false, res.Error
	}
}

// 根据IDS获取菜单信息
func (a *MenuRepository) GetByIds(ids []int64) (result []auth.Menu) {
	sql := `
			SELECT 
				a.id, a.fid, a.title, a.name, TRIM(a.icon) as icon, a.name as path, a.node_level, a.component, a.out_page,
				IFNULL((SELECT 1 FROM auth_menu b WHERE  b.fid=a.id  LIMIT 1),0) as has_sub_node
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

	//1.菜单可能被分配给  tb_auth_casbin_rules 的权限
	sql := `
		DELETE FROM auth_casbin_rule WHERE fr_auth_post_mount_has_menu_button_id IN (
			SELECT 
				CONCAT(organization_id, ':', menu_id, ':', button_id) 
			FROM 
				auth_organization_menu_button 
			WHERE 
				menu_id = ?
			)
		)
	`
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("MenuRepository 删除 auth_casbin_rule 失败", res.Error.Error())
	}

	//2. 菜单可能被分配给组织机构的权限关联数据
	sql = `
		DELETE FROM 
			auth_organization_menu_button  
		WHERE 
			menu_id = ?
	`
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("MenuRepository 删除 auth_organization_menu_button 失败", res.Error.Error())
	}

	//3. 菜单可能被分配给组织机构的权限按钮数据
	sql = `DELETE FROM auth_organization_menu WHERE menu_id=?`
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("MenuRepository 删除 auth_organization_menu 失败", res.Error.Error())
	}

	//4. 删除菜单关联的待分配按钮子表
	sql = `DELETE FROM auth_menu_button WHERE menu_id  = ? `
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("MenuRepository 删除 auth_menu_button 失败", res.Error.Error())
	}
}
