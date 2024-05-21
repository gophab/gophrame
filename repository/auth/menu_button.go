package auth

import (
	"strconv"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"

	"github.com/wjshen/gophrame/domain/auth"

	"gorm.io/gorm"
)

type MenuButtonRepository struct {
	*gorm.DB `inject:"database"`
}

var menuButtonRepository *MenuButtonRepository = &MenuButtonRepository{}

func init() {
	inject.InjectValue("menuButtonRepository", menuButtonRepository)
}

func (a *MenuButtonRepository) getCounts(sysMenuId int64) (count int64) {
	if res := a.Model(&auth.MenuButton{}).Where("menu_id=?", sysMenuId).Count(&count); res.Error == nil {
		return count
	}
	return 0
}

// 查询
func (a *MenuButtonRepository) List(sysMenuId int64) (counts int64, data []auth.MenuButton) {
	counts = a.getCounts(sysMenuId)
	if counts > 0 {
		sql := `
			SELECT
				a.menu_id,
				a.button_id,
				b.cn_name AS button_name,
				a.request_method,  
				a.request_url,
				a.status,
				a.remark
			FROM
				auth_menu_button a 
			LEFT JOIN auth_button b ON a.button_id=b.id
			WHERE   
				a.menu_id=?
		`
		if res := a.Raw(sql, sysMenuId).Find(&data); res.Error != nil {
			logger.Error("AuthSystemMenuButtonModel 查询出错：", res.Error.Error())
		}
	} else {
		return 0, nil
	}
	return
}

// 新增
func (a *MenuButtonRepository) InsertData(list auth.ButtonArray) bool {

	//注意: 这里必须使用  Table 函数指定表名
	// 不能使用  a.Model(a) 设置表名，a.Model 函数会设置a上绑定的很多结构信息,这样就会导致与外部的数据类型 ButtonArray 不一致，最终 gorm 反射出错
	if res := a.Table((&auth.MenuButton{}).TableName()).Create(&list); res.Error != nil {
		logger.Error("系统菜单子表插入数据出错：", res.Error.Error())
		return false
	}
	return true
}

// 更新
func (a *MenuButtonRepository) UpdateData(menuEdit auth.MenuEdit) bool {

	// 删除可能存在的垃圾数据(由于开发测试阶段可能存储手动添加的测试、忘记删除最终会导致权限出现异常)
	var list = menuEdit.ButtonArray

	var newIds string
	for _, item := range list {
		if item.ButtonId > 0 {
			newIds += strconv.FormatInt(item.ButtonId, 10) + ","
		}
	}

	sql := `DELETE FROM auth_menu_button WHERE menu_id=? AND !FIND_IN_SET(button_id,?)`
	if res := a.Exec(sql, menuEdit.Id, newIds); res.Error != nil {
		logger.Error("删除可能的垃圾数据出错：", res.Error.Error())
		return false
	}

	for index, item := range list {
		// 如果 id 为 0 表示修改的过程中新增了数据
		sql := `
			INSERT INTO auth_menu_button (menu_id, button_id, request_url, request_method, status, remark)
			SELECT ?,?,?,?,1,? FROM DUAL WHERE NOT EXISTS(
				SELECT 1 FROM auth_menu_button a WHERE a.menu_id=? AND a.button_id=?
			)
		`
		if res := a.Exec(sql, item.MenuId, item.ButtonId, item.RequestUrl, item.RequestMethod, item.Remark, item.MenuId, item.ButtonId); res.Error != nil {
			logger.Error("修改界面ID未生成，新增了数据，执行sql出错：", res.Error.Error())
			return false
		} else if res.RowsAffected <= 0 {
			sql = `
				UPDATE 
					auth_menu_button  
				SET  
					request_url=?,
					request_method=?,
					remark=?
				WHERE 
					menu_id=? 
					AND button_id=?
			`
			if res = a.Exec(sql, item.RequestUrl, item.RequestMethod, item.Remark, item.MenuId, item.ButtonId); res.Error != nil {
				logger.Error("修改界面ID未生成，针对已输入文本继续做了修改，执行sql出错：", res.Error.Error())
				return false
			}
		}
		list[index] = item
	}
	return true
}

// 删除数据
func (a *MenuButtonRepository) DeleteData(id int64) bool {
	if res := a.Delete(a, id); res.Error == nil {
		return true
	} else {
		logger.Error("MenuButtonRepository 数据删除失败", res.Error.Error())
	}
	return false
}

// 批量删除数据
func (a *MenuButtonRepository) BatchDeleteData(ids []struct{ menu_id, button_id int64 }) bool {
	values := []string{}
	for _, v := range ids {
		values = append(values, strconv.FormatInt(v.menu_id, 10)+":"+strconv.FormatInt(v.button_id, 10))
	}
	sql := `DELETE FROM auth_menu_button WHERE FIND_IN_SET(CONCAT(menu_id, ':', button_id),?)`
	if res := a.Exec(sql, values); res.Error == nil {
		return true
	} else {
		logger.Error("MenuButtonRepository 批量数据删除失败", res.Error.Error())
	}
	return false
}

// 新增
func (a *MenuButtonRepository) InsertMap(data map[string]interface{}) bool {
	a.Model(&auth.MenuButton{}).Create(&data)
	return true
}

// 根据菜单ID获取按钮信息
func (a *MenuButtonRepository) MenuButton(menuId int64) (data []auth.MenuButton) {
	sql := `
		SELECT 
			a.menu_id, 
			a.button_id,
			a.remark, 
			a.request_url, 
			a.request_method, 
			a.status, 
			b.cn_name as button_name 
		FROM 
			auth_menu_button a 
		LEFT JOIN 
			auth_button b ON a.button_id = b.id 
		WHERE 
			a.menu_id = ?
	`
	a.Raw(sql, menuId).Scan(&data)
	return
}

// 数据更新hook函数，负责更新菜单被引用的地方，同步更新
func (a *MenuButtonRepository) UpdateHook(menuId int64) {
	// 更新菜单挂接的按钮之后，可能存在按钮被删除，因此需要删除的数据主要有：1. tb_auth_casbin_rules 表被应用的按钮数据
	sql := `
		DELETE FROM 
			auth_casbin_rule 
		WHERE 
			fr_auth_post_mount_has_menu_button_id IN (
				SELECT 
					CONCAT(b.organization_id, ':', b.menu_id, ':', b.button_id)
				FROM 
					auth_organization_menu a, 
					auth_organization_menu_button b
				WHERE 
					b.organization_id = a.organization_id
					AND b.menu_id = a.menu_id
					AND a.menu_id=?
					AND b.button_id NOT IN (
						SELECT 
							d.button_id 
						FROM  
							auth_menu c, 
							auth_menu_button  d   
						WHERE 
							c.id=d.menu_id
							AND c.id=?
					)
			)	
	`
	if res := a.Exec(sql, menuId, menuId); res.Error != nil {
		logger.Error("MenuButtonRepository UpdateHook 删除 auth_casbin_rule 关联按钮数据出错", res.Error.Error())
	}

	sql = `
		DELETE FROM 
			auth_organization_menu_button 
		WHERE 
			menu_id=?	
	`
	if res := a.Exec(sql, menuId, menuId); res.Error != nil {
		logger.Error("MenuButtonRepository UpdateHook 删 auth_organization_menu_button 关联按钮数据出错", res.Error.Error())
	}

	// 批量更新菜单被引用的所有地方
	sql = `
		UPDATE (
			SELECT DISTINCT 
				a.id AS menu_id,
				b.button_id,
				b.request_method,
				b.request_url,
				c.organization_id
			FROM  
				auth_menu a, 
				auth_menu_button b, 
				auth_organization_menu c, 
				auth_organization_menu_button d
			WHERE 
				a.id=b.menu_id
				AND c.menu_id = b.menu_id 
				AND d.organization_id=c.organization_id
				AND d.menu_id=b.menu_id
				AND d.button_id=b.button_id
				AND a.id=? 
		)  AS f 
		LEFT JOIN auth_casbin_rule e 
			ON e.fr_auth_post_mount_has_menu_button_id=CONCAT(f.organization_id,':',f.menu_id,':',f.button_id)
		SET 
			e.v1=f.request_url, 
			e.v2=f.request_method 
		WHERE  
			e.ptype='p'  
			AND LENGTH(IFNULL(f.request_url,''))>0 
			AND LENGTH(IFNULL(f.request_method,''))>0 
		`
	if res := a.Exec(sql, menuId); res.Error != nil {
		logger.Error("MenuButtonRepository UpdateHook 更新 tb_auth_casbin_rule 出错", res.Error.Error())
	}
}

// 判断按钮是否系统菜单引用
func (a *MenuButtonRepository) GetByButtonId(buttonId int64) bool {
	data := []auth.MenuButton{}
	a.Model(&auth.MenuButton{}).Where("button_id = ?", buttonId).Find(&data)
	return len(data) == 0
}

// GetSystemMenuButtonList 待分配的系统菜单、按钮 数据列表
// 注意：按钮的id有可能和主菜单id重复，所以按钮id基准值增加 100000 （10万），后续分配权限时减去 10万即可
func (a *MenuButtonRepository) GetSystemAuthorities() (int64, []auth.MenuButton) {
	data := []auth.MenuButton{}
	sql := `
		SELECT
			a.id AS id,
			a.fid AS fid,
			a.title As title,
			'menu'  AS node_type,
			(CASE WHEN a.fid=0 THEN 1 ELSE 0 END) AS expand,
			a.sort
		FROM
			auth_menu a
		UNION  
		SELECT 
			IFNULL(c.id,0)+? AS id,
			IFNULL(b.menu_id,0) AS fid,
			IFNULL(c.cn_name,'') AS title,
			'button' AS node_type,
			0 AS expand,
			0 AS sort
		FROM
			auth_menu_button b   
		LEFT JOIN 
			auth_button c ON b.button_id=c.id
		ORDER BY 
			sort DESC, fid ASC, id ASC
	`
	if res := a.Raw(sql, 100000).Find(&data); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, data
	} else {
		logger.Error("查询系统待分配菜单出错：", res.Error.Error())
	}

	return 0, data
}
