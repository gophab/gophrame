package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/util"

	"github.com/wjshen/gophrame/default/domain/auth"

	"gorm.io/gorm"
)

type AuthorityRepository struct {
	*gorm.DB `inject:"database"`
}

var authorityRepository *AuthorityRepository = &AuthorityRepository{}

func init() {
	inject.InjectValue("authorityRepository", authorityRepository)
}

// 查询用户在指定页面拥有的按钮列表
func (u *AuthorityRepository) GetButtonListByMenuId(roleIds []string, MenuId int64) (r []auth.Button) {
	sql := `
		SELECT  
			c.*
		FROM  
			auth_role_menu a,
			auth_role_menu_button b,
			auth_button c 
		WHERE
			a.role_id=b.role_id   
		AND
			a.menu_id=b.menu_id
		AND
			b.button_id=c.id
		AND 
			a.role_id IN  (?)
		AND
			a.menu_id = ?
		`
	if res := u.Raw(sql, roleIds, MenuId).Find(&r); res.Error != nil {
		logger.Error("获取指定页面(菜单)所拥有的按钮权限出错", res.Error.Error())
	}
	return
}

// GetSystemMenuButtonList 待分配的系统菜单、按钮 数据列表
// 注意：按钮的id有可能和主菜单id重复，所以按钮id基准值增加 100000 （10万），后续分配权限时减去 10万即可
func (a *AuthorityRepository) GetSystemMenuButtonList() (counts int64, data []auth.AuthNode) {
	var menuNodes []auth.AuthNode
	sql := `
		SELECT 
			a.id AS id,
			a.fid AS fid,
			a.title,
			'menu' AS node_type,
			(CASE WHEN a.fid=0 THEN 1 ELSE 0 END) AS expand,
			a.sort
		FROM
			auth_menu a
		ORDER BY sort DESC, fid ASC, id ASC
	`
	if err := a.Raw(sql).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []auth.AuthNode{}
	}

	var buttonNodes []auth.AuthNode
	sql = `  
		SELECT 
			IFNULL(c.id,0) AS id,
			IFNULL(b.menu_id,0) AS fid,
			IFNULL(c.cn_name,'') AS title,
			'button' AS node_type,
			0 AS expand,
			0 AS sort
		FROM
			auth_menu_button b   
		LEFT JOIN auth_button c ON b.button_id=c.id
		ORDER BY sort DESC, fid ASC, id ASC
	`
	if err := a.Raw(sql).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(&menuNodes, &buttonNodes)
	}

	if err := util.CreateSqlResFormatFactory().ScanToTreeData(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []auth.AuthNode{}
	}
}

// 已分配给部门、岗位的系统菜单、按钮
func (a *AuthorityRepository) GetAssignedMenuButtonList(roleId string) (counts int64, data []auth.AuthNode) {
	var menuNodes []auth.AuthNode
	sql := `
		SELECT  
			b.id 										AS id,
			b.fid 									AS fid, 
			b.title									AS title,
			'menu' 									AS node_type,
			(case when b.fid=0 then 1 else 0 end) AS expand,
			b.sort 									AS sort
		FROM 
			auth_role_menu a, auth_menu b  
		WHERE  
			a.menu_id=b.id
			AND a.status=1
			AND a.role_id=?
		ORDER BY sort DESC, id ASC, fid ASC
	`
	if err := a.Raw(sql, roleId).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []auth.AuthNode{}
	}

	var buttonNodes []auth.AuthNode
	sql = `
		SELECT  
			IFNULL(d.id,0)				AS id,
			c.menu_id 						AS fid,
			IFNULL(d.cn_name,'') 	AS title,
			'button' 							AS node_type,
			0 										AS expand, 
			d.id 									AS sort
		FROM 
			auth_role_menu a, 
			auth_role_menu_button c, 
			auth_button d  
		WHERE
			a.role_id=c.role_id
			AND a.menu_id=c.menu_id
			AND c.button_id=d.id
			AND a.status=1
			AND a.role_id=?
		ORDER BY sort DESC, id ASC, fid ASC
	`
	if err := a.Raw(sql, roleId).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(&menuNodes, &buttonNodes)
	}

	if err := util.CreateSqlResFormatFactory().ScanToTreeData(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []auth.AuthNode{}
	}
}

// 给组织机构（部门、岗位）分配菜单权限
func (a *AuthorityRepository) AssginAuthForOrg(orgId, menuId, buttonId int64, nodeType string) (assignRes bool) {
	assignRes = true
	sql := `INSERT  INTO auth_role_menu(role_id, menu_id)
					SELECT ?,? FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM auth_role_menu a WHERE a.role_id=? AND a.menu_id=? FOR UPDATE)
			`
	//1.菜单分配权限
	if nodeType == "menu" {
		// 每一个菜单都可能有上级菜单，最大支持到三级菜单即可
		var tmpMenuId = menuId
		var tmpFid int64 = 0
		for i := 0; i < 3; i++ {
			if res := a.Exec(sql, orgId, menuId, orgId, menuId); res.Error == nil {
				tmpSql := "select a.fid from auth_menu a where a.id=?"
				if _ = a.Raw(tmpSql, tmpMenuId).First(&tmpFid); tmpFid > 0 {
					tmpMenuId = tmpFid
				}
			} else {
				assignRes = false
				logger.Error("auth_post_mount_has_menu  插入 menuList 时出错", res.Error.Error())
			}
		}
	}

	//2.按钮权限分配
	if nodeType == "button" {
		if buttonId > 0 {
			sql = `
				INSERT INTO auth_role_menu_button(role_id, menu_id, button_id)
				SELECT ?,?,? FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM auth_role_menu_button a WHERE a.role_id=? AND a.menu_id=? AND a.button_id=? FOR UPDATE)
			`
			if res := a.Exec(sql, orgId, menuId, buttonId, orgId, menuId, buttonId); res.Error == nil {
				// 3.继续分配接口的访问权限(casbin_rules写入相关数据)
				assignRes = a.AssginCasbinAuthPolicyToOrg(fmt.Sprintf("%d:%d:%d", orgId, menuId, buttonId), nodeType)
			} else {
				logger.Error("auth_post_mount_has_menu_button  表分配按钮失败", res.Error.Error())
				assignRes = false
			}
		}
	}
	return assignRes
}

// 从组织机构（部门、岗位）删除权限
func (a *AuthorityRepository) DeleteAuthFromOrg(roleId, menuId, buttonId int, nodeType string) bool {
	if nodeType == "menu" {
		sql := "DELETE   FROM auth_role_menu WHERE role_id=? AND menu_id=?"
		if res := a.Exec(sql, roleId, menuId); res.Error == nil {
			return true
		}
	} else if nodeType == "button" {
		sql := "DELETE FROM auth_role_menu_button WHERE role_id=? AND menu_id=? AND button_id=?"
		if res := a.Exec(sql, roleId, menuId, buttonId); res.Error == nil {
			return a.DeleteCasbibRules(fmt.Sprintf("%d:%d:%d", roleId, menuId, buttonId), nodeType)
		}
	}
	return false
}

// 删除 casbin 表接口已分配的权限
func (a *AuthorityRepository) DeleteCasbibRules(authPostMountHasMenuButtonId string, nodeType string) (resBool bool) {
	resBool = true
	if nodeType == "button" {
		sql := "DELETE FROM auth_casbin_rule WHERE fr_auth_post_mount_has_menu_button_id=? AND ptype='p' "
		if res := a.Exec(sql, authPostMountHasMenuButtonId); res.Error != nil {
			// 角色继承关系暂时不删除，只要删除相关的节点权限即可
			logger.Error("AuthMenuAssignRepository 删除casbin权限失败", res.Error.Error())
			resBool = false
		}
	}
	return
}

// 给组织机构节点分配casbin的policy策略权限
func (a *AuthorityRepository) AssginCasbinAuthPolicyToOrg(authPostMountHasMenuButtonId string, nodeType string) (resBool bool) {
	// 参见 69 行注释
	var failTryTimes = 1
	resBool = true
	// 分配了按钮，才可以同步分配按钮对应的路由接口
	if nodeType == "button" {
		segs := strings.Split(authPostMountHasMenuButtonId, ":")
		// 首先给组织机构分配p权限(polic权限)
		sql := `
			SELECT   
				'p' as ptype,
				b.role_id,
				c.request_url,
				UPPER(c.request_method) AS request_method
			FROM  
				auth_role_menu_button a, 
				auth_role_menu b, 
				auth_menu_button c
			WHERE   
				a.role_id=b.role_id
				AND a.menu_id = b.menu_id 
				AND b.menu_id = c.menu_id 
				AND c.button_id = a.button_id
				AND a.role_id = ?
				AND a.menu_id = ?
				AND a.button_id = ?
		`
		var tmp struct {
			Ptype         string
			RoleId        int
			RequestUrl    string
			RequestMethod string
		}
		if res := a.Raw(sql, segs[0], segs[1], segs[2]).First(&tmp); res.Error == nil && tmp.Ptype != "" {
			sql = `
			INSERT  INTO auth_casbin_rule(ptype,v0,v1,v2,fr_auth_post_mount_has_menu_button_id,v3,v4,v5)
			SELECT  ?,?,?,?,?,'','',''  FROM   DUAL 
			WHERE NOT EXISTS(SELECT 1 FROM auth_casbin_rule a force index(idx_vp01) WHERE a.ptype=? AND a.v0=? AND a.v1=? AND a.v2=? FOR UPDATE)
			`
		label1:
			if res = a.Exec(sql, tmp.Ptype, tmp.RoleId, tmp.RequestUrl, tmp.RequestMethod, authPostMountHasMenuButtonId, tmp.Ptype, tmp.RoleId, tmp.RequestUrl, tmp.RequestMethod); res.Error == nil {
				// 为当前节点继续分配g(group权限，设置部门继承关系)
				return a.AssginCasbinAuthGroupToOrg(tmp.RoleId)
			} else {
				if failTryTimes <= 5 {
					failTryTimes++
					goto label1
				}
				resBool = false
				logger.Error("AuthMenuAssignRepository 发生错误", res.Error.Error())
			}
		} else {
			resBool = false
			logger.Error("根据参数：authPostMountHasMenuButtonId 查询时出错：", "authPostMountHasMenuButtonId", authPostMountHasMenuButtonId, res.Error.Error())
		}
	}
	return resBool
}

// 给组织机构节点分配casbin的group（角色继承关系权限）
func (a *AuthorityRepository) AssginCasbinAuthGroupToOrg(orgId int) (resBool bool) {
	// 参见 69 行注释
	var failTryTimes = 1
	resBool = true
	sql := "SELECT path_info FROM sys_role WHERE id =?"
	var pathInfo string
	if res := a.Raw(sql, orgId).First(&pathInfo); res.Error == nil {
		if len(pathInfo) > 0 {
			orgIdArray := strings.Split(pathInfo, ",")
			orgLen := len(orgIdArray)
			sql = `
				INSERT INTO auth_casbin_rule (ptype,v0,v1,v2,v3,v4,v5) 
				SELECT 'g',?,?,'','','',''  FROM   DUAL   
				WHERE NOT EXISTS(SELECT 1 FROM auth_casbin_rule a WHERE a.ptype='g' AND v0=? AND v1=? FOR UPDATE)
				`
			var lastId = 0
			var id = 0
			var err error
			for i := 1; i <= orgLen; i++ {
				// 遍历组织机构id
				if id, err = strconv.Atoi(orgIdArray[orgLen-i]); err == nil && i > 1 && id > 0 {
				label:
					if res = a.Exec(sql, lastId, id, lastId, id); res.Error != nil {
						if failTryTimes <= 5 {
							failTryTimes++
							goto label
						}
						logger.Error("AuthMenuAssignRepository 批量插入角色继承关系时出错", res.Error.Error())
						resBool = false
					}
				}
				lastId = id
			}
		}
	} else {
		resBool = false
	}
	return resBool
}

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetUserAuthorities(userId string) (result []auth.AuthNode) {
	sql := `
		SELECT GROUP_CONCAT(b.path_info) id
		FROM
			sys_role_user a
		LEFT JOIN 
			sys_role b ON a.role_id = b.id
		WHERE 
			a.status=1 AND a.user_id = ?
		GROUP BY a.user_id
	`
	var orgPathInfo string
	if res := a.Raw(sql, userId).First(&orgPathInfo); res.Error != nil {
		logger.Error("查询用户所在的岗位、部门、公司节点出错：" + res.Error.Error())
		return result
	} else if len(orgPathInfo) == 0 {
		return result
	}

	// L1
	var deptNodes []auth.AuthNode
	sql = `
		SELECT   
			c.id,
			c.fid AS fid,
			c.name AS title, 
			'dept' 	AS node_type,
			1 AS expand
		FROM 
			sys_role c
		WHERE 
			FIND_IN_SET(c.id,?) 
			AND c.status=1
	`
	a.Raw(sql, orgPathInfo).Scan(&deptNodes)
	if len(deptNodes) == 0 {
		return
	}

	// L2
	var menuNodes []auth.AuthNode
	sql = `  
		SELECT DISTINCT  
			e.id, 
			(CASE 
				WHEN e.fid=0 THEN d.role_id 
				ELSE e.fid 
			END) AS fid,
			e.title,
			'menu' AS node_type,
			(CASE WHEN e.fid=0 THEN 1 ELSE 0 END) AS expand
		FROM
			auth_role_menu d, 
			auth_menu  e 
		WHERE  
			FIND_IN_SET(d.role_id,?)
			AND d.menu_id=e.id
	`
	a.Raw(sql, orgPathInfo).Scan(&menuNodes)

	if len(menuNodes) > 0 {
		// L3
		var buttonNodes []auth.AuthNode
		sql = `
			SELECT
				g.id AS id,
				f.menu_id AS fid ,
				g.cn_name AS title,
				'button' AS  node_type,
				0 AS expand
			FROM  
				auth_role_menu d ,
				auth_role_menu_button f,
				auth_button  g
			WHERE  
				d.role_id=f.role_id
				AND d.menu_id=f.menu_id
				AND f.button_id=g.id
				AND d.status=1 
				AND f.status=1 
				AND FIND_IN_SET(d.role_id,?)
		`
		a.Raw(sql, orgPathInfo).Scan(&buttonNodes)
		if len(buttonNodes) > 0 {
			a.makeChildren(&menuNodes, &buttonNodes)
		}

		var menuTree []auth.AuthNode
		if err := util.CreateSqlResFormatFactory().ScanToTreeData(menuNodes, &menuTree); err != nil {
			logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
			return
		}

		a.makeChildren(&deptNodes, &menuTree)
	}

	if err := util.CreateSqlResFormatFactory().ScanToTreeData(deptNodes, &result); err != nil {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
	}

	return
}

func (a *AuthorityRepository) makeChildren(fnodes, cnodes *[]auth.AuthNode) {
	var fMap = map[int64]*auth.AuthNode{}
	for i, n := range *fnodes {
		fMap[n.Id] = &(*fnodes)[i]
	}
	for _, n := range *cnodes {
		fnode := fMap[n.Fid]
		if fnode != nil {
			fnode.Children = append(fnode.Children, n)
		}
	}
}
