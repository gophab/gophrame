package repository

import (
	"github.com/gophab/gophrame/module/slink/domain"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/transaction"

	"gorm.io/gorm"
)

type ShortLinkRepository struct {
	*gorm.DB `inject:"database"`
}

var shortLinkReposistory = &ShortLinkRepository{}

func init() {
	inject.InjectValue("shortLinkRepository", shortLinkReposistory)
}

func (r *ShortLinkRepository) GetById(id string) (*domain.ShortLink, error) {
	var result domain.ShortLink
	if res := transaction.Session().Model(&domain.ShortLink{}).Where("id=?", id).Find(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *ShortLinkRepository) DeleteById(id string) error {
	return transaction.Session().Delete(&domain.ShortLink{}, "id=?", id).Error
}

func (r *ShortLinkRepository) GetByKey(key string) (*domain.ShortLink, error) {
	var result domain.ShortLink
	if res := transaction.Session().Model(&domain.ShortLink{}).Where("`key`=?", key).Find(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r *ShortLinkRepository) DeleteByKey(key string) error {
	return transaction.Session().Delete(&domain.ShortLink{}, "`key`=?", key).Error
}

func (r *ShortLinkRepository) CreateShortLink(shortLink *domain.ShortLink) (*domain.ShortLink, error) {
	if res := transaction.Session().Create(shortLink); res.Error == nil && res.RowsAffected > 0 {
		return shortLink, nil
	} else {
		return nil, res.Error
	}
}

func (r *ShortLinkRepository) ExpireExpiredShortLinks() {
	if res := r.Delete(&domain.ShortLink{}, "expired_time < CURRENT_TIMESTAMP"); res.Error != nil {
		logger.Warn("Expire shortlink error: ", res.Error.Error())
	}
}
