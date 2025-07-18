package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/util/collection"

	"github.com/gophab/gophrame/module/authority/domain"

	OperationModel "github.com/gophab/gophrame/module/operation/domain"

	"gorm.io/gorm"
)

type AuthorityRepository struct {
	*gorm.DB `inject:"database"`
}

var authorityRepository *AuthorityRepository = &AuthorityRepository{}

func init() {
	inject.InjectValue("authorityRepository", authorityRepository)
}

func (a *AuthorityRepository) GetRoleAuthority(roleId string, authType string, authId string) (*domain.RoleAuthority, error) {
	var result *domain.RoleAuthority
	if res := a.Model(&domain.RoleAuthority{}).Where("role_id = ? and auth_type = ? and auth_id = ?", roleId, authType, authId).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (a *AuthorityRepository) GetRoleAuthorities(roleId string, authType string) ([]*domain.RoleAuthority, error) {
	var results = make([]*domain.RoleAuthority, 0)
	if res := a.Model(&domain.RoleAuthority{}).Where("role_id = ? and auth_type = ?", roleId, authType).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return results, nil
	} else {
		return []*domain.RoleAuthority{}, res.Error
	}
}

func (a *AuthorityRepository) GetRolesAuthorities(roleIds []string, authType string) ([]*domain.RoleAuthority, error) {
	var results = make([]*domain.RoleAuthority, 0)
	if res := a.Model(&domain.RoleAuthority{}).Where("role_id in ? and auth_type = ?", roleIds, authType).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return results, nil
	} else {
		return []*domain.RoleAuthority{}, res.Error
	}
}

func (a *AuthorityRepository) GetUserAuthority(userId string, authType string, authId string) (*domain.UserAuthority, error) {
	var result *domain.UserAuthority
	if res := a.Model(&domain.UserAuthority{}).Where("user_id = ? and auth_type = ? and auth_id = ?", userId, authType, authId).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (a *AuthorityRepository) GetUserAuthorities(userId string, authType string) ([]*domain.UserAuthority, error) {
	var results = make([]*domain.UserAuthority, 0)
	if res := a.Model(&domain.UserAuthority{}).Where("user_id = ? and auth_type = ?", userId, authType).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return results, nil
	} else {
		return []*domain.UserAuthority{}, res.Error
	}
}

func (a *AuthorityRepository) GetOrganizationAuthority(organizationId string, authType string, authId string) (*domain.OrganizationAuthority, error) {
	var result *domain.OrganizationAuthority
	if res := a.Model(&domain.OrganizationAuthority{}).Where("organization_id = ? and auth_type = ? and auth_id = ?", organizationId, authType, authId).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return result, nil
	} else {
		return nil, res.Error
	}
}

func (a *AuthorityRepository) GetOrganizationAuthorities(organizationId string, authType string) ([]*domain.OrganizationAuthority, error) {
	var results = make([]*domain.OrganizationAuthority, 0)
	if res := a.Model(&domain.OrganizationAuthority{}).Where("organization_id = ? and auth_type = ?", organizationId, authType).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return results, nil
	} else {
		return []*domain.OrganizationAuthority{}, res.Error
	}
}

func (a *AuthorityRepository) GetOrganizationsAuthorities(organizationIds []string, authType string) ([]*domain.OrganizationAuthority, error) {
	var results = make([]*domain.OrganizationAuthority, 0)
	if res := a.Model(&domain.OrganizationAuthority{}).Where("organization_id in ? and auth_type = ?", organizationIds, authType).Find(&results); res.Error == nil && res.RowsAffected > 0 {
		return results, nil
	} else {
		return []*domain.OrganizationAuthority{}, res.Error
	}
}

// 给角色授权
func (a *AuthorityRepository) SetAuthoritiesByRoleId(roleId string, authType string, authIds []string) {
	// 1. clear
	// TODO: DeleteRoleOperations
	a.DeleteRoleAuthoritiesByAuthType(roleId, authType)

	// 2. add
	for _, authId := range authIds {
		var authority = &domain.RoleAuthority{
			RoleId: roleId,
			Authority: domain.Authority{
				AuthType: authType,
				AuthId:   authId,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist RoleAuthority error: ", res.Error.Error())
			break
		}
	}
}

// 给用户授权
func (a *AuthorityRepository) SetAuthoritiesByUserId(userId string, authType string, authIds []string) {
	// 1. clear
	// TODO: DeleteRoleOperations
	a.DeleteUserAuthoritiesByAuthType(userId, authType)

	// 2. add
	for _, authId := range authIds {
		var authority = &domain.UserAuthority{
			UserId: userId,
			Authority: domain.Authority{
				AuthType: authType,
				AuthId:   authId,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist UserAuthority error: ", res.Error.Error())
			break
		}
	}
}

// 给用户授权
func (a *AuthorityRepository) SetAuthoritiesByOrganizationId(organizationId string, authType string, authIds []string) {
	// 1. clear
	// TODO: DeleteRoleOperations
	a.DeleteOrganizationAuthoritiesByAuthType(organizationId, authType)

	// 2. add
	for _, authId := range authIds {
		var authority = &domain.OrganizationAuthority{
			OrganizationId: organizationId,
			Authority: domain.Authority{
				AuthType: authType,
				AuthId:   authId,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist OrganizationAuthority error: ", res.Error.Error())
			break
		}
	}
}

func (a *AuthorityRepository) DeleteRoleAuthority(roleId string, authType string, authId string) {
	var authority = &domain.RoleAuthority{
		RoleId: roleId,
		Authority: domain.Authority{
			AuthType: authType,
			AuthId:   authId,
		},
	}
	if res := a.Delete(&authority); res.Error != nil {
		logger.Warn("Delete UserAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteRoleAuthorities(roleId string) {
	if res := a.Delete(&domain.RoleAuthority{}, "role_id = ?", roleId); res.Error != nil {
		logger.Warn("Delete RoleAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteRoleAuthoritiesByAuthType(roleId string, authType string) {
	if res := a.Delete(&domain.RoleAuthority{}, "role_id = ? AND auth_type = ?", roleId, authType); res.Error != nil {
		logger.Warn("Delete RoleAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteUserAuthority(userId string, authType string, authId string) {
	var authority = &domain.UserAuthority{
		UserId: userId,
		Authority: domain.Authority{
			AuthType: authType,
			AuthId:   authId,
		},
	}
	if res := a.Delete(&authority); res.Error != nil {
		logger.Warn("Delete UserAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteUserAuthoritiesByAuthType(userId string, authType string) {
	if res := a.Delete(&domain.UserAuthority{}, "user_id = ? AND auth_type = ?", userId, authType); res.Error != nil {
		logger.Warn("Delete UserAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteUserAuthorities(userId string) {
	if res := a.Delete(&domain.UserAuthority{}, "user_id = ?", userId); res.Error != nil {
		logger.Warn("Delete UserAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteOrganizationAuthority(organizationId string, authType string, authId string) {
	var authority = &domain.OrganizationAuthority{
		OrganizationId: organizationId,
		Authority: domain.Authority{
			AuthType: authType,
			AuthId:   authId,
		},
	}
	if res := a.Delete(&authority); res.Error != nil {
		logger.Warn("Delete OrganizationAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteOrganizationAuthorities(organizationId string) {
	if res := a.Delete(&domain.OrganizationAuthority{}, "organization_id = ?", organizationId); res.Error != nil {
		logger.Warn("Delete OrganizationAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteOrganizationAuthoritiesByAuthType(organizationId string, authType string) {
	if res := a.Delete(&domain.OrganizationAuthority{}, "organization_id = ? AND auth_type = ?", organizationId, authType); res.Error != nil {
		logger.Warn("Delete OrganizationAuthority error: ", res.Error.Error())
	}
}

// ////////////////////////////////////////
// Special for Operation (Menu & Button)
// ////////////////////////////////////////
func (u *AuthorityRepository) GetMenuByRoleIds(roleIds []string) (result []*OperationModel.Menu, err error) {
	sql := `
		SELECT  
			distinct b.*
		FROM 
			auth_role_authority a, auth_menu b  
		WHERE  
			a.role_id IN ? AND a.auth_type='menu' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY b.sort ASC, b.id ASC, b.fid ASC
	`
	if err = u.Raw(sql, roleIds).Find(&result).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	return
}

// 查询用户在指定页面拥有的按钮列表
func (u *AuthorityRepository) GetButtonListByMenuId(roleIds []string, menuId int64) (r []*OperationModel.Button, err error) {
	sql := `
		SELECT  
			distinct c.*
		FROM  
			auth_role_authority a,
			auth_button c 
		WHERE
			a.auth_type='button'   
			AND a.auth_id=c.id
			AND	c.fid = ?
		`
	if len(roleIds) > 0 {
		sql += ` AND a.role_id IN  (?) `
	}
	sql += `
		ORDER BY c.sort ASC, c.fid ASC, c.id ASC
	`
	if err = u.Raw(sql, menuId, roleIds).Find(&r).Error; err != nil {
		logger.Error("获取指定页面(菜单)所拥有的按钮权限出错", err.Error())
	}
	return
}

// GetSystemOperations 待分配的系统菜单、按钮 数据列表
// 注意：按钮的id有可能和主菜单id重复，所以按钮id基准值增加 100000 （10万），后续分配权限时减去 10万即可
func (a *AuthorityRepository) GetSystemOperations() (counts int64, data []*OperationModel.Operation) {
	var menuNodes []*OperationModel.Operation
	sql := `
		SELECT 
			a.id AS id,
			a.fid AS fid,
			a.title AS title,
			a.name AS name,
			a.tags AS tags,
			'menu' AS type,
			(CASE WHEN a.fid=0 THEN 1 ELSE 0 END) AS expand,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=a.id) AS leaf,
			(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=a.id) AS cks,
			a.sort
		FROM
			auth_menu a
		ORDER BY sort ASC, fid ASC, id ASC
	`
	if err := a.Raw(sql).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []*OperationModel.Operation{}
	}

	var buttonNodes []*OperationModel.Operation
	sql = `  
		SELECT 
			b.id AS id,
			b.fid AS fid,
			b.title AS title,
			b.name AS name,
			b.tags AS tags,
			'button' AS type,
			0 AS expand,
			0 AS sort
		FROM
			auth_button b   
		ORDER BY fid ASC, id ASC
	`
	if err := a.Raw(sql).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(menuNodes, buttonNodes)
	}

	if err := a.makeTree(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []*OperationModel.Operation{}
	}
}

func (a *AuthorityRepository) makeTree(src []*OperationModel.Operation, dest *[]*OperationModel.Operation) error {
	var result = *dest
	var srcMap = make(map[string]*OperationModel.Operation)
	for _, item := range src {
		srcMap[item.Id] = item
	}
	for _, item := range src {
		if item.Fid != "" && item.Fid != "0" {
			var parent = srcMap[item.Fid]
			if parent != nil {
				if parent.Children == nil {
					parent.Children = make([]*OperationModel.Operation, 0)
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

// 已分配给部门、岗位的系统菜单、按钮
func (a *AuthorityRepository) GetOperationsByRoleId(roleId string) (counts int64, data []*OperationModel.Operation) {
	var menuNodes []*OperationModel.Operation
	sql := `
		SELECT  
			b.id 										AS id,
			b.fid 									AS fid, 
			b.title									AS title,
			'menu' 									AS type,
			(case when b.fid=0 then 1 else 0 end) AS expand,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=b.id) AS leaf,
			(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=b.id) AS cks,
			b.sort 									AS sort
		FROM 
			auth_role_authority a, auth_menu b  
		WHERE  
			a.role_id=? && a.auth_type='menu' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY sort DESC, id ASC, fid ASC
	`
	if err := a.Raw(sql, roleId).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []*OperationModel.Operation{}
	}

	var buttonNodes []*OperationModel.Operation
	sql = `
		SELECT  
			b.id									AS id,
			b.fid 								AS fid,
			b.title 							AS title,
			'button' 							AS type,
			0 										AS expand, 
			b.id 									AS sort
		FROM 
			auth_role_authority a, auth_button b  
		WHERE
			a.role_id=? && a.auth_type='button' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY b.id ASC, b.fid ASC
	`
	if err := a.Raw(sql, roleId).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(menuNodes, buttonNodes)
	}

	if err := a.makeTree(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []*OperationModel.Operation{}
	}
}

func (a *AuthorityRepository) GetOperationsByRoleIds(roleIds []string) (result []*OperationModel.Operation) {
	var authorities = make(map[string]*domain.Authority)

	var tmp = make([]*domain.RoleAuthority, 0)
	if a.Where("role_id IN ?", roleIds).Find(&tmp).Error == nil && len(tmp) > 0 {
		for _, auth := range tmp {
			authorities[fmt.Sprintf("%s:%s", auth.AuthType, auth.AuthId)] = &auth.Authority
		}
	}

	// 4. 合并
	if len(authorities) > 0 {
		var menuIds = make([]string, 0)
		for _, auth := range authorities {
			if auth.AuthType == "menu" {
				menuIds = append(menuIds, auth.AuthId)
			}
		}

		if len(menuIds) > 0 {
			var menuNodes []*OperationModel.Operation = make([]*OperationModel.Operation, 0)
			sql := `  
				SELECT
					e.id AS id, 
					e.fid AS fid,
					e.name AS name,
					e.title AS title,
					'menu' AS type,
					(CASE WHEN e.fid=0 THEN 1 ELSE 0 END) AS expand,
					(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=e.id) AS leaf,
					(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=e.id) AS cks,
					e.sort AS sort
				FROM
					auth_menu e 
				WHERE  
					e.id IN ?
				ORDER BY sort ASC, fid ASC
			`
			a.Raw(sql, menuIds).Scan(&menuNodes)

			if len(menuNodes) > 0 {
				var buttonIds = make([]string, 0)
				for _, auth := range authorities {
					if auth.AuthType == "button" {
						buttonIds = append(buttonIds, auth.AuthId)
					}
				}
				if len(buttonIds) > 0 {
					var buttonNodes []*OperationModel.Operation = make([]*OperationModel.Operation, 0)
					sql = `
						SELECT
							g.id AS id,
							g.fid AS fid,
							g.name AS name,
							g.name AS title,
							0 AS sort,
							'button' AS  type,
							0 AS expand
						FROM  
							auth_button  g
						WHERE  
							g.id IN ?
						ORDER BY sort ASC, fid ASC
					`
					a.Raw(sql, buttonIds).Scan(&buttonNodes)
					if len(buttonNodes) > 0 {
						a.makeChildren(menuNodes, buttonNodes)
					}
				}
			}

			if err := a.makeTree(menuNodes, &result); err != nil {
				logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
				return
			}
		}
	}
	return
}

// 给角色分配系统菜单、按钮
func (a *AuthorityRepository) SetOperationsByRoleId(roleId string, data []*OperationModel.Operation) {
	// 1. clear
	// TODO: DeleteRoleOperations
	a.DeleteRoleOperations(roleId)

	// 2. add
	for _, operation := range data {
		var authority = &domain.RoleAuthority{
			RoleId: roleId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist RoleAuthority error: ", res.Error.Error())
			break
		}
	}
}

func (a *AuthorityRepository) AddOperationsByRoleId(roleId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		var authority = &domain.RoleAuthority{
			RoleId: roleId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist RoleAuthority error: ", res.Error.Error())
			break
		}
	}
}

// 给角色分配系统菜单、按钮
func (a *AuthorityRepository) SetOpertionsByUserId(userId string, data []*OperationModel.Operation) {
	// 1. clear
	a.DeleteUserAuthorities(userId)

	// 2. add
	for _, operation := range data {
		var authority = &domain.UserAuthority{
			UserId: userId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist UserAuthority error: ", res.Error.Error())
			break
		}
	}
}

func (a *AuthorityRepository) DeleteOperationsByRoleId(roleId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		a.DeleteRoleAuthority(roleId, operation.Type, operation.Id)
	}
}

func (a *AuthorityRepository) DeleteOpertionsByUserId(userId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		a.DeleteUserAuthority(userId, operation.Type, operation.Id)
	}
}

func (a *AuthorityRepository) AddOpertionsByUserId(userId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		var authority = &domain.UserAuthority{
			UserId: userId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist UserAuthority error: ", res.Error.Error())
			break
		}
	}
}

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetOpertionsByUserId(userId string) (count int64, data []*OperationModel.Operation) {
	var menuNodes []*OperationModel.Operation
	sql := `
		SELECT  
			b.id 										AS id,
			b.fid 									AS fid, 
			b.title									AS title,
			'menu' 									AS type,
			(case when b.fid=0 then 1 else 0 end) AS expand,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=b.id) AS leaf,
			(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=b.id) AS cks,
			b.sort 									AS sort
		FROM 
			auth_user_authority a, auth_menu b  
		WHERE  
			a.user_id=? && a.auth_type='menu' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY sort DESC, id ASC, fid ASC
	`
	if err := a.Raw(sql, userId).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []*OperationModel.Operation{}
	}

	var buttonNodes []*OperationModel.Operation
	sql = `
		SELECT  
			b.id									AS id,
			b.fid 								AS fid,
			b.title 							AS title,
			'button' 							AS type,
			0 										AS expand, 
			b.id 									AS sort
		FROM 
			auth_user_authority a, auth_button b  
		WHERE
			a.user_id=? && a.auth_type='button' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY b.id ASC, b.fid ASC
	`
	if err := a.Raw(sql, userId).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(menuNodes, buttonNodes)
	}

	if err := a.makeTree(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []*OperationModel.Operation{}
	}
}

// 给角色分配系统菜单、按钮
func (a *AuthorityRepository) SetOperationsByOrganizationId(organizationId string, data []*OperationModel.Operation) {
	// 1. clear
	a.DeleteOrganizationOperations(organizationId)

	// 2. add
	for _, operation := range data {
		var authority = &domain.OrganizationAuthority{
			OrganizationId: organizationId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist OrganizationAuthority error: ", res.Error.Error())
			break
		}
	}
}

func (a *AuthorityRepository) AddOperationsByOrganizationId(organizationId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		var authority = &domain.OrganizationAuthority{
			OrganizationId: organizationId,
			Authority: domain.Authority{
				AuthType: operation.Type,
				AuthId:   operation.Id,
				Status:   1,
			},
		}
		if res := a.FirstOrCreate(authority); res.Error != nil {
			logger.Warn("Persist OrganizationAuthority error: ", res.Error.Error())
			break
		}
	}
}

func (a *AuthorityRepository) DeleteOperationsByOrganizationId(organizationId string, data []*OperationModel.Operation) {
	for _, operation := range data {
		a.DeleteOrganizationAuthority(organizationId, operation.Type, operation.Id)
	}
}

func (a *AuthorityRepository) DeleteRoleOperations(roleId string) {
	if res := a.Delete(&domain.RoleAuthority{}, "role_id = ? and auth_type in ?", roleId, []string{"menu", "button"}); res.Error != nil {
		logger.Warn("Delete RoleAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteUserOperations(userId string) {
	if res := a.Delete(&domain.UserAuthority{}, "user_id = ? and auth_type in ?", userId, []string{"menu", "button"}); res.Error != nil {
		logger.Warn("Delete UserAuthority error: ", res.Error.Error())
	}
}

func (a *AuthorityRepository) DeleteOrganizationOperations(organizationId string) {
	if res := a.Delete(&domain.OrganizationAuthority{}, "organization_id = ? and auth_type in ?", organizationId, []string{"menu", "button"}); res.Error != nil {
		logger.Warn("Delete OrganizationAuthority error: ", res.Error.Error())
	}
}

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetOperationsByOrganizationId(organizationId string) (count int64, data []*OperationModel.Operation) {
	var menuNodes []*OperationModel.Operation
	sql := `
		SELECT  
			b.id 										AS id,
			b.fid 									AS fid, 
			b.title									AS title,
			'menu' 									AS type,
			(case when b.fid=0 then 1 else 0 end) AS expand,
			(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=b.id) AS leaf,
			(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=b.id) AS cks,
			b.sort 									AS sort
		FROM 
			auth_organization_authority a, auth_menu b  
		WHERE  
			a.organization_id=? && a.auth_type='menu' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY sort DESC, id ASC, fid ASC
	`
	if err := a.Raw(sql, organizationId).Find(&menuNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
		return
	}

	if len(menuNodes) == 0 {
		return 0, []*OperationModel.Operation{}
	}

	var buttonNodes []*OperationModel.Operation
	sql = `
		SELECT  
			b.id									AS id,
			b.fid 								AS fid,
			b.title 							AS title,
			'button' 							AS type,
			0 										AS expand, 
			b.id 									AS sort
		FROM 
			auth_organization_authority a, auth_button b  
		WHERE
			a.organization_id=? && a.auth_type='button' AND a.auth_id=b.id
			AND a.status=1
		ORDER BY b.id ASC, b.fid ASC
	`
	if err := a.Raw(sql, organizationId).Find(&buttonNodes).Error; err != nil {
		logger.Error("查询系统待分配菜单出错：", err.Error())
	}

	if len(buttonNodes) > 0 {
		a.makeChildren(menuNodes, buttonNodes)
	}

	if err := a.makeTree(menuNodes, &data); err == nil {
		return int64(len(data)), data
	} else {
		logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
		return 0, []*OperationModel.Operation{}
	}
}

func (a *AuthorityRepository) GetOperationsByOrganizationIds(organizationIds []string) (result []*OperationModel.Operation) {
	var authorities = make(map[string]*domain.Authority)

	// 3. organization -> organization_authorities
	var tmp = make([]*domain.OrganizationAuthority, 0)
	if a.Where("organization_id IN ?", organizationIds).Find(&tmp).Error == nil && len(tmp) > 0 {
		for _, auth := range tmp {
			authorities[fmt.Sprintf("%s:%s", auth.AuthType, auth.AuthId)] = &auth.Authority
		}
	}

	// 4. 合并
	if len(authorities) > 0 {
		var menuIds = make([]string, 0)
		for _, auth := range authorities {
			if auth.AuthType == "menu" {
				menuIds = append(menuIds, auth.AuthId)
			}
		}

		if len(menuIds) > 0 {
			var menuNodes []*OperationModel.Operation = make([]*OperationModel.Operation, 0)
			sql := `  
				SELECT
					e.id AS id, 
					e.fid AS fid,
					e.title AS title,
					e.sort AS sort,
					'menu' AS type,
					(CASE WHEN e.fid=0 THEN 1 ELSE 0 END) AS expand,
					(SELECT CASE WHEN COUNT(*) =0 THEN 1 ELSE 0 END FROM auth_button WHERE fid=e.id) AS leaf,
					(SELECT GROUP_CONCAT(name) FROM auth_button WHERE fid=e.id) AS cks,
				FROM
					auth_role_authority d, 
					auth_menu e 
				WHERE  
					d.auth_type = 'menu'
					AND FIND_IN_SET(d.auth_id, ?)
					AND d.auth_id=e.id
			`
			a.Raw(sql, menuIds).Scan(&menuNodes)

			if len(menuNodes) > 0 {
				var buttonIds = make([]string, 0)
				for _, auth := range authorities {
					if auth.AuthType == "button" {
						buttonIds = append(buttonIds, auth.AuthId)
					}
				}
				if len(buttonIds) > 0 {
					var buttonNodes []*OperationModel.Operation = make([]*OperationModel.Operation, 0)
					sql = `
						SELECT
							g.id AS id,
							g.fid AS fid ,
							g.title AS title,
							0 AS sort,
							'button' AS  type,
							0 AS expand
						FROM  
							auth_role_authority d,
							auth_button  g
						WHERE  
							d.auth_type = 'button'
							AND FIND_IN_SET(d.auth_id, ?)
							AND d.auth_id=g.id
					`
					a.Raw(sql, buttonIds).Scan(&buttonNodes)
					if len(buttonNodes) > 0 {
						a.makeChildren(menuNodes, buttonNodes)
					}
				}
			}

			if err := a.makeTree(menuNodes, &result); err != nil {
				logger.Error("AuthSystemMenuModel 树形化出错:" + err.Error())
				return
			}
		}
	}
	return
}

// func (a *AuthorityRepository) loadRoleIds(roleIds collection.Set[string], roleId string) {
// 	if roleIds.Has(roleId) {
// 		return
// 	}
// 	roleIds.Add(roleId)

// 	var role domain.Role
// 	if res := a.Model(&domain.Role{}).Where("id = ?", roleId).First(&role); res.Error == nil && res.RowsAffected > 0 {
// 		if role.Includes != "" {
// 			var includeIds = strings.Split(role.Includes, ",")
// 			for _, rid := range includeIds {
// 				a.loadRoleIds(roleIds, rid)
// 			}
// 		}
// 	}
// }

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetRoleOperations(roleId string) (result []*OperationModel.Operation) {
	var roleIdSet = make(collection.Set[string])
	// a.loadRoleIds(roleIdSet, roleId)
	roleIdSet.Add(roleId)
	return a.GetOperationsByRoleIds(roleIdSet.AsList())
}

func (a *AuthorityRepository) GetRolesOperations(roleIds []string) (result []*OperationModel.Operation) {
	var roleIdSet = make(collection.Set[string])
	roleIdSet.AddAll(roleIds)
	return a.GetOperationsByRoleIds(roleIdSet.AsList())
}

// func (a *AuthorityRepository) loadOrganizationIds(organizationIds collection.Set[string], organizationId string) {
// 	if organizationIds.Has(organizationId) {
// 		return
// 	}
// 	organizationIds.Add(organizationId)

// 	var organization domain.Organization
// 	if res := a.Model(&domain.Organization{}).Where("id = ?", organizationId).First(&organization); res.Error == nil && res.RowsAffected > 0 {
// 		if organization.PathInfo != "" {
// 			var oids = strings.Split(organization.PathInfo, ",")
// 			organizationIds.AddAll(oids)
// 		}
// 	}
// }

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetOrganizationOperations(organizationId string) (result []*OperationModel.Operation) {
	var organizationIdSet = make(collection.Set[string])
	organizationIdSet.Add(organizationId)
	// a.loadOrganizationIds(organizationIdSet, organizationId)
	return a.GetOperationsByRoleIds(organizationIdSet.AsList())
}

func (a *AuthorityRepository) GetOrganizationsOperations(organizationIds []string) (result []*OperationModel.Operation) {
	var organizationIdSet = make(collection.Set[string])
	// for _, organizationId := range organizationIds {
	// 	a.loadOrganizationIds(organizationIdSet, organizationId)
	// }
	organizationIdSet.AddAll(organizationIds)
	return a.GetOperationsByOrganizationIds(organizationIdSet.AsList())
}

// 根据用户id查询已经分配的菜单
func (a *AuthorityRepository) GetUserAvailableOperations(userId string) (result []*OperationModel.Operation) {
	var operations = make(map[string]*OperationModel.Operation)

	// 1. user -> organization -> organization_authorities
	sql := `
		SELECT 
			organization_id
		FROM
			sys_organization_user
		WHERE
			user_id = ?
	`
	var organizationIds = make([]string, 0)
	if res := a.Raw(sql, userId).Find(&organizationIds); res.Error == nil && len(organizationIds) > 0 {
		var ops = a.GetOrganizationsOperations(organizationIds)
		if len(ops) > 0 {
			for _, op := range ops {
				operations[fmt.Sprintf("%s:%s", op.Type, op.Id)] = op
			}
		}
	}

	// 2. user -> roles -> role_authorities
	sql = `
		SELECT 
			role_id
		FROM
			sys_role_user
		WHERE
			user_id = ?
	`
	var roleIds = make([]string, 0)
	if res := a.Raw(sql, userId).Find(&roleIds); res.Error == nil && len(roleIds) > 0 {
		var ops = a.GetRolesOperations(roleIds)
		if len(ops) > 0 {
			for _, op := range ops {
				operations[fmt.Sprintf("%s:%s", op.Type, op.Id)] = op
			}
		}
	}

	// 3. user -> user_authorities
	_, ops := a.GetOpertionsByUserId(userId)
	if len(ops) > 0 {
		for _, op := range ops {
			operations[fmt.Sprintf("%s:%s", op.Type, op.Id)] = op
		}
	}

	if len(operations) > 0 {
		result = make([]*OperationModel.Operation, 0)
		for _, op := range operations {
			result = append(result, op)
		}
	}
	return
}

// 删除 casbin 表接口已分配的权限
func (a *AuthorityRepository) DeleteCasbinRules(authPostMountHasMenuButtonId string, nodeType string) (resBool bool) {
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

func (a *AuthorityRepository) makeChildren(fnodes, cnodes []*OperationModel.Operation) {
	var fMap = map[string]*OperationModel.Operation{}
	for _, n := range fnodes {
		fMap[n.Id] = n
	}
	for _, n := range cnodes {
		fnode := fMap[n.Fid]
		if fnode != nil {
			fnode.Children = append(fnode.Children, n)
		}
	}
}
