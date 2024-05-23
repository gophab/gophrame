package repository

import (
	"github.com/wjshen/gophrame/core/inject"

	"github.com/wjshen/gophrame/default/domain"

	"gorm.io/gorm"
)

type UserOptionRepository struct {
	*gorm.DB `inject:"database"`
}

var userOptionRepository = &UserOptionRepository{}

func init() {
	inject.InjectValue("userOptionRepository", userOptionRepository)
}

func (r *UserOptionRepository) GetUserOptions(userId string) (*domain.UserOptions, error) {
	result := &domain.UserOptions{UserId: userId, Options: make(map[string]domain.UserOption)}

	var userOptions []domain.UserOption
	if res := r.Where("user_id=?", userId).Find(&userOptions); res.Error == nil && res.RowsAffected > 0 {
		for _, option := range userOptions {
			result.Options[option.Name] = option
		}
		return result, nil
	} else {
		return nil, res.Error
	}

}

func (r *UserOptionRepository) RemoveAllUserOptions(userId string) error {
	return r.Delete(&domain.UserOption{UserId: userId}).Error
}
