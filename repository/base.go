package repository

import (
	"reflect"

	"github.com/gophab/gophrame/core/transaction"
	"gorm.io/gorm"
)

type BaseRepository[T any, K any] struct {
	*gorm.DB `inject:"database"`
}

func (r BaseRepository[T, K]) GetById(id K) (*T, error) {
	var result T
	if res := r.Where("id = ?", id).First(&result); res.Error == nil && res.RowsAffected > 0 {
		return &result, nil
	} else {
		return nil, res.Error
	}
}

func (r BaseRepository[T, K]) SaveAll(list []*T) (result []*T, err error) {
	transaction.Session().Begin()
	defer func() {
		err = transaction.Session().Commit().Error
	}()

	for _, item := range list {
		transaction.Session().Save(item)
	}

	result = list

	return
}

func (r BaseRepository[T, K]) DeleteById(id K) error {
	var result T
	if _, b := reflect.TypeOf(result).FieldByName("DelFlag"); b {
		res := r.Model(&result).Where("id = ?", id).Update("del_flag", true)
		return res.Error
	} else {
		res := r.Where("id = ?", id).Delete(&result)
		return res.Error
	}
}
