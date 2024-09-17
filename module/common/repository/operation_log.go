package repository

import (
	"fmt"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/common/domain"

	"gorm.io/gorm"
)

type OperationLogRepository struct {
	*gorm.DB `inject:"database"`
}

var operationLogRepository = &OperationLogRepository{}

func init() {
	inject.InjectValue("operationLogRepository", operationLogRepository)
}

func (r *OperationLogRepository) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.OperationLog, error) {
	tx := r.Model(&domain.OperationLog{})

	for k, v := range conds {
		tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	tx.Order("operation_time desc")

	var results = []*domain.OperationLog{}
	if res := query.Page(tx, pageable).Find(&results); res.Error == nil {
		return -1, results, nil
	} else {
		return 0, results, res.Error
	}
}

func (r *OperationLogRepository) Append(log *domain.OperationLog) error {
	return r.Create(log).Error
}
