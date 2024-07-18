package transaction

import (
	"context"
	"time"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/routine"
	"gorm.io/gorm"
)

type DBSession struct {
	*gorm.DB
	inTransaction bool
}

type TransactionManager struct {
	*gorm.DB `inject:"database"`
	session  *routine.ThreadLocal[*DBSession]
}

var transactionManager = &TransactionManager{session: routine.NewThreadLocal(&DBSession{})}

func init() {
	inject.InjectValue("transactionManager", transactionManager)
}

func (m *TransactionManager) InTransaction() bool {
	dbSession := m.session.Get()
	if dbSession != nil {
		return dbSession.inTransaction
	}
	return false
}

func (m *TransactionManager) BeginTransaction() {
	if !m.InTransaction() {
		m.session.Set(&DBSession{
			DB:            m.Session().Begin(),
			inTransaction: true,
		})
	}
}

func (m *TransactionManager) Commit() error {
	if m.InTransaction() {
		err := m.session.Get().Commit().Error
		m.session.Remove()
		return err
	}
	return nil
}

func (m *TransactionManager) Rollback() {
	if m.InTransaction() {
		m.session.Get().Rollback()
		m.session.Remove()
	}
}

func (m *TransactionManager) EndTransaction() error {
	if m.InTransaction() {
		defer m.session.Remove()
		if err := m.session.Get().Commit().Error; err != nil {
			m.session.Get().Rollback()
			return err
		}
	}
	return nil
}

func (m *TransactionManager) Session() *gorm.DB {
	if m.session.Get() == nil {
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
		m.session.Set(&DBSession{
			DB:            m.DB.WithContext(timeoutContext),
			inTransaction: false,
		})
	}
	return m.session.Get().DB
}

func Begin() {
	transactionManager.BeginTransaction()
}

func Session() *gorm.DB {
	return transactionManager.Session()
}

func Commit() error {
	return transactionManager.Commit()
}

func Rollback() {
	transactionManager.Rollback()
}

func End() error {
	return transactionManager.EndTransaction()
}
