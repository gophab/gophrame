package repository

import (
	"errors"
	"time"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/system/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	*gorm.DB `inject:"database"`
}

var userRepository *UserRepository = &UserRepository{}

func init() {
	logger.Info("Initializing User Repository")
	inject.InjectValue("userRepository", userRepository)
}

func (h *UserRepository) CheckUser(username, password string) (bool, error) {
	var user domain.User
	if res := h.Select("id").
		Where("login=? OR mobile=? OR email=?", username, username, username).
		Where("password=?", util.SHA1(password)).
		Where("del_flag=?", false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) GetUserByUserNamePassword(username, password string) (*domain.User, error) {
	var user domain.User
	if res := h.Select("id").
		Where("login=? OR mobile=? OR email=?", username, username, username).
		Where("password=?", util.SHA1(password)).
		Where("del_flag=?", false).
		First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return &user, nil
}

func (h *UserRepository) CheckUserLogin(username string) (bool, error) {
	var user domain.User
	if res := h.Where("login = ? AND del_flag = ?", username, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) CheckUserMobile(username string) (bool, error) {
	var user domain.User
	if res := h.Where("mobile = ? AND del_flag = ?", username, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) CheckUserEmail(username string) (bool, error) {
	var user domain.User
	if res := h.Where("email = ? AND del_flag = ? ", username, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) CheckUserLoginId(login string, id string) (bool, error) {
	var user domain.User
	if res := h.Where("login = ? AND id != ? AND del_flag = ? ", login, id, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) CheckUserMobileId(mobile string, id string) (bool, error) {
	var user domain.User
	if res := h.Where("mobile = ? AND id != ? AND del_flag = ?", mobile, id, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) CheckUserEmailId(email string, id string) (bool, error) {
	var user domain.User
	if res := h.Where("email = ? AND id != ? AND del_flag = ?", email, id, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) ExistUserByID(id string) (bool, error) {
	var user domain.User
	if res := h.Select("id").Where("id = ? AND del_flag = ?", id, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return false, res.Error
	}

	return true, nil
}

func (h *UserRepository) GetUserTotal(maps any) (int64, error) {
	var count int64
	if err := h.Model(&domain.User{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) buildQuery(conds map[string]any) *gorm.DB {
	q := r.Model(&domain.User{})
	for k, v := range conds {
		switch k {
		case "search":
			q = q.Where("name like ? or mobile like ? or telephone like ? or email like ? or id = ?", "%"+v.(string)+"%", "%"+v.(string)+"%", "%"+v.(string)+"%", "%"+v.(string)+"%", v)
		case "name":
			q = q.Where("name like ?", "%"+v.(string)+"%")
		case "mobile":
			q = q.Where("mobile like ?", "%"+v.(string)+"%")
		case "login":
			q = q.Where("login like ?", "%"+v.(string)+"%")
		case "email":
			q = q.Where("email like ?", "%"+v.(string)+"%")
		case "ids":
			q = q.Where("id in ?", v)
		default:
			q = q.Where(k+"=?", v)
		}
	}
	q = q.Where("del_flag = ?", false)
	return q
}

func (r *UserRepository) GetAll(conds map[string]any) ([]*domain.User, error) {
	var users []*domain.User = make([]*domain.User, 0)

	var q = r.buildQuery(conds)
	if res := q.Preload("Roles").Find(&users); res.Error == nil {
		return users, nil
	} else {
		return []*domain.User{}, res.Error
	}
}

func (r *UserRepository) Find(conds map[string]any, pageable query.Pageable) (total int64, list []*domain.User) {
	var tx = r.buildQuery(conds)
	total = 0
	if !pageable.NoCount() {
		if tx.Count(&total).Error != nil || total == 0 {
			return
		}
	}

	query.Page(tx, pageable).Find(&list)
	return
}

func (h *UserRepository) GetUser(username string) (*domain.User, error) {
	var user domain.User
	err := h.Preload("Roles").Where("(login = ? OR mobile = ? OR email = ?) AND del_flag = ? ", username, username, username, false).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, nil
}

func (h *UserRepository) GetUserByLogin(login string) (*domain.User, error) {
	var user domain.User
	if res := h.Preload("Roles").Where("login = ? AND del_flag = ? ", login, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return &user, nil
}

func (h *UserRepository) GetUserByMobile(mobile string) (*domain.User, error) {
	var user domain.User
	if res := h.Preload("Roles").Where("mobile = ? AND del_flag = ? ", mobile, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return &user, nil
}

func (h *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if res := h.Preload("Roles").Where("email = ? AND del_flag = ? ", email, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return &user, nil
}

func (h *UserRepository) GetUserById(id string) (*domain.User, error) {
	var user domain.User
	if res := h.Preload("Roles").Where("id = ? AND del_flag = ? ", id, false).First(&user); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return &user, nil
}

func (h *UserRepository) GetUserByIds(ids []string) ([]*domain.User, error) {
	var users []*domain.User
	if res := h.Preload("Roles").Where("id in ? AND del_flag = ? ", ids, false).Find(&users); res.Error != nil || res.RowsAffected <= 0 {
		return nil, res.Error
	}

	return users, nil
}

func (h *UserRepository) UpdateUser(entity *domain.User) error {
	var user domain.User
	if res := h.Where("id = ? AND del_flag = ? ", entity.Id, false).Find(&user); res.Error != nil {
		return res.Error
	} else if res.RowsAffected <= 0 {
		return errors.New("user not found")
	}

	if err := util.CopyFieldsExcept(&user, *entity, "LastLoginTime", "LastLoginIp", "Password", "CreatedTime", "CreatedBy"); err != nil {
		return err
	}

	if entity.Password != "####*****####" && entity.Password != "" {
		user.SetPassword(entity.Password)
	}

	// roles
	var roles []domain.Role
	if len(entity.Roles) > 0 {
		ids := []string{}
		for _, v := range entity.Roles {
			ids = append(ids, v.Id)
		}
		h.Where("id in (?)", ids).Find(&roles)
	}
	h.Model(&user).Association("Roles").Replace(roles)

	// columns
	h.Model(&user).Omit("created_by", "created_time", "last_login_time", "last_login_ip", "login_times").Save(&user)

	return nil
}

func (r *UserRepository) Change(userId string, column string, value any) (int64, error) {
	if res := r.Model(&domain.User{}).Where("id = ?", userId).UpdateColumn(column, value); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, nil
	} else {
		return 0, res.Error
	}
}

func (r *UserRepository) Changes(userId string, columns map[string]any) (int64, error) {
	if res := r.Model(&domain.User{}).Where("id = ?", userId).UpdateColumns(util.DbFields(columns)); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, nil
	} else {
		return 0, res.Error
	}
}

func (r *UserRepository) ConditionChange(conds map[string]any, column string, value any) (int64, error) {
	q := r.buildQuery(conds)
	if res := q.UpdateColumn(column, value); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, nil
	} else {
		return 0, res.Error
	}
}

func (r *UserRepository) ConditionChanges(conds map[string]any, columns map[string]any) (int64, error) {
	q := r.buildQuery(conds)
	if res := q.UpdateColumns(util.DbFields(columns)); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, nil
	} else {
		return 0, res.Error
	}
}

func (h *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	if user.Password != "####*****####" {
		user.SetPassword(user.Password)
	}

	var roles []domain.Role
	if len(user.Roles) > 0 {
		ids := []string{}
		for _, v := range user.Roles {
			ids = append(ids, v.Id)
		}
		h.Where("id in (?)", ids).Find(&roles)
	}
	if err := h.Create(&user).Association("Roles").Append(roles); err != nil {
		return nil, err
	}
	return user, nil
}

func (h *UserRepository) DeleteUser(id string) error {
	var user domain.User
	if res := h.Where("id = ? AND del_flag = ?", id, false).Find(&user); res == nil || res.RowsAffected <= 0 {
		return res.Error
	}

	// 删除相关角色
	h.Model(&user).Association("Roles").Delete()

	// 删除对象
	if err := h.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (h *UserRepository) CleanAllUser() error {
	if err := h.Unscoped().Where("del_flag = ?", false).Delete(&domain.User{}).Error; err != nil {
		return err
	}

	return nil
}

func (h *UserRepository) GetUsersAll() ([]*domain.User, error) {
	var users []*domain.User
	err := h.Where("del_flag = ?", false).Preload("Roles").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 根据关键词查询用户表的条数
func (u *UserRepository) getCounts(userName string) (counts int64) {
	sql := "select count(*) as counts from sys_user WHERE (login like ? or mobile like ? or email like ? or name like ?) AND del_flag = ?"
	if _ = u.Raw(sql, "%"+userName+"%", "%"+userName+"%", "%"+userName+"%", "%"+userName+"%", false).First(&counts); counts > 0 {
		return counts
	} else {
		return 0
	}
}

// 权限分配查询（包含用户岗位信息）
func (a *UserRepository) GetUserWithOrganizations(userName string, pageable query.Pageable) (totalCounts int64, list []*domain.User) {
	totalCounts = a.getCounts(userName)
	if totalCounts > 0 {
		sql := `
			SELECT  
				a.*, 
				(
					SELECT  
						REPLACE(IFNULL(GROUP_CONCAT(name ORDER BY id ASC),''),',',' | ')
					FROM sys_organization b
					WHERE 
						b.id IN (SELECT organization_id FROM sys_organization_user c WHERE c.user_id=a.id AND c.status=1)
				) organization_name 
			FROM 
				sys_user a 
			WHERE 
				(login LIKE ? OR mobile LIKE ? OR email LIKE ? OR name LIKE ?) 
				AND del_flag = ?
			LIMIT ?,?
		`
		if res := a.Raw(sql, "%"+userName+"%", "%"+userName+"%", "%"+userName+"%", "%"+userName+"%", false, pageable.GetOffset(), pageable.GetLimit()).Find(&list); res.RowsAffected > 0 {
			return totalCounts, list
		} else {
			return totalCounts, nil
		}
	}

	return 0, nil
}

func (a *UserRepository) LogUserLogin(userId string, loginIp string) error {
	sql := `UPDATE sys_user SET login_times = login_times + 1, last_login_time=?, last_login_ip=? WHERE id=?`
	return a.Exec(sql, time.Now(), loginIp, userId).Error
}
