package repository

import (
	"github.com/wjshen/gophrame/core/inject"

	"github.com/wjshen/gophrame/default/domain"

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
