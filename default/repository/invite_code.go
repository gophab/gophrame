package repository

import (
	"time"

	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/default/domain"

	"gorm.io/gorm"
)

type InviteCodeRepository struct {
	*gorm.DB `inject:"database"`
}

var inviteCodeRepository = &InviteCodeRepository{}

func init() {
	inject.InjectValue("inviteCodeRepository", inviteCodeRepository)
}

func (s *InviteCodeRepository) FindByInviteCode(inviteCode string) (*domain.InviteCode, error) {
	var result domain.InviteCode
	if res := s.Where("invite_code=?", inviteCode).Where("del_flag=?", false).First(&result); res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected <= 0 || result.IsExpired() {
		return nil, nil
	}

	return &result, nil
}

func (s *InviteCodeRepository) GetUserInviteCode(userId string, channel string) (*domain.InviteCode, error) {
	var result domain.InviteCode
	if res := s.Where("user_id=?", userId).Where("channel=?", channel).Where("del_flag=?", false).Where("expire_time is NULL or expire_time > ?", time.Now()).First(&result); res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected <= 0 || result.IsExpired() {
		return nil, nil
	}
	return &result, nil
}
