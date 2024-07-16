package transaction

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/routine"
	"gorm.io/gorm"
)

type TransactionManager struct {
	*gorm.DB `inject:"database"`
	session  *routine.ThreadLocal[*gorm.DB]
}

var transactionManager = &TransactionManager{session: routine.NewThreadLocal(&gorm.DB{})}

func init() {
	inject.InjectValue("transactionManager", transactionManager)
}

func (m *TransactionManager) InTransaction() bool {
	return m.session.Get() != nil
}

func (m *TransactionManager) BeginTransaction() {
	if !m.InTransaction() {
		// m.session.Set(m.Begin())
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
	if m.InTransaction() {
		return m.session.Get()
	}
	return m.DB
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
