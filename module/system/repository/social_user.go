package repository

import (
	"time"

	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/system/domain"

	"gorm.io/gorm"
)

type SocialUserRepository struct {
	*gorm.DB `inject:"database"`
}

var socialUserRepository *SocialUserRepository = &SocialUserRepository{}

func init() {
	inject.InjectValue("socialUserRepository", socialUserRepository)
}

func (r *SocialUserRepository) GetById(id string) (*domain.SocialUser, error) {
	var result domain.SocialUser
	if res := r.Where("id=?", id).Where("del_flag=?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SocialUserRepository) GetBySocialId(socialType string, socialId string) (*domain.SocialUser, error) {
	var result domain.SocialUser
	if res := r.Where("type=?", socialType).Where("social_id=?", socialId).Where("del_flag=?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SocialUserRepository) GetByUserId(socialType string, userId string) (*domain.SocialUser, error) {
	var result domain.SocialUser
	if res := r.Where("type=?", socialType).Where("user_id=?", userId).Where("del_flag=?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SocialUserRepository) GetByMobile(socialType string, mobile string) (*domain.SocialUser, error) {
	var result domain.SocialUser
	if res := r.Where("type=?", socialType).Where("mobile=?", mobile).Where("del_flag=?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SocialUserRepository) GetByEmail(socialType string, email string) (*domain.SocialUser, error) {
	var result domain.SocialUser
	if res := r.Where("type=?", socialType).Where("email=?", email).Where("del_flag=?", false).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *SocialUserRepository) buildQuery(conds map[string]any) *gorm.DB {
	q := r.Model(&domain.SocialUser{})
	for k, v := range conds {
		switch k {
		case "search":
			q = q.Where("name like ? or description like ? or model like ?", "%"+v.(string)+"%", "%"+v.(string)+"%", "%"+v.(string)+"%")
		case "name":
			q = q.Where("name like ?", "%"+v.(string)+"%")
		case "description":
			q = q.Where("description like ?", "%"+v.(string)+"%")
		case "location":
			q = q.Where("location like ?", "%"+v.(string)+"%")
		case "start_time": // 已开始
			q = q.Where("start_time <= ?", v)
		case "end_time": // 已截止
			q = q.Where("end_time <= ?", v)
		case "ids":
			q = q.Where("id in ?", v)
		default:
			q = q.Where(k+"=?", v)
		}
	}
	q = q.Where("del_flag = ?", false)

	return q
}

func (r *SocialUserRepository) Change(id string, column string, value any) error {
	return r.Changes(id, core.M{column: value})
}

func (r *SocialUserRepository) Changes(id string, columns map[string]any) error {
	if _, b := columns["last_modified_by"]; !b {
		columns["last_modified_by"] = SecurityUtil.GetCurrentUserId(nil)
	}
	if _, b := columns["last_modified_time"]; !b {
		columns["last_modified_time"] = time.Now()
	}

	return r.Model(&domain.SocialUser{}).Where("id = ?", id).UpdateColumns(util.DbFields(columns)).Error
}

func (r *SocialUserRepository) ConditionChange(conds map[string]any, column string, value any) (int64, error) {
	return r.ConditionChanges(conds, core.M{column: value})
}

func (r *SocialUserRepository) ConditionChanges(conds map[string]any, columns map[string]any) (int64, error) {
	if _, b := columns["last_modified_by"]; !b {
		columns["last_modified_by"] = SecurityUtil.GetCurrentUserId(nil)
	}
	if _, b := columns["last_modified_time"]; !b {
		columns["last_modified_time"] = time.Now()
	}

	q := r.buildQuery(util.DbFields(conds))
	if res := q.UpdateColumns(util.DbFields(columns)); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, nil
	} else {
		return 0, res.Error
	}
}
