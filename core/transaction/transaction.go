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
	cancel        context.CancelFunc
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
	dbSession := m._session()
	if dbSession != nil {
		return dbSession.inTransaction
	}
	return false
}

func (m *TransactionManager) BeginTransaction() {
	var session = m._session()
	defer func() {
		if session != nil && session.cancel != nil {
			session.cancel()
		}
		m.session.Remove()
	}()

	if !session.inTransaction {
		m.session.Set(&DBSession{
			DB:            session.DB.Begin(),
			inTransaction: true,
			cancel:        session.cancel,
		})
	}
}

func (m *TransactionManager) Commit() error {
	var session = m._session()
	defer func() {
		if session != nil && session.cancel != nil {
			session.cancel()
		}
		m.session.Remove()
	}()

	if session.inTransaction {
		return session.Commit().Error
	}
	return nil
}

func (m *TransactionManager) Rollback() {
	var session = m._session()
	defer func() {
		if session != nil && session.cancel != nil {
			session.cancel()
		}
		m.session.Remove()
	}()
	if session.inTransaction {
		session.Rollback()
	}
}

func (m *TransactionManager) EndTransaction() error {
	var session = m._session()
	defer func() {
		if session != nil && session.cancel != nil {
			session.cancel()
		}
		m.session.Remove()
	}()
	if session.inTransaction {
		if err := session.Commit().Error; err != nil {
			session.Rollback()
			return err
		}
	}
	return nil
}

func (m *TransactionManager) _session() *DBSession {
	if m.session.Get() == nil {
		timeoutContext, cancel := context.WithTimeout(context.Background(), 1200*time.Second)
		m.session.Set(&DBSession{
			DB:            m.DB.WithContext(timeoutContext),
			inTransaction: false,
			cancel:        cancel,
		})
	}
	return m.session.Get()
}

func (m *TransactionManager) Session() *gorm.DB {
	result := m._session().DB
	// if !result.DryRun {
	// 	m.session.Remove()
	// 	result = m._session().DB
	// }
	return result
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
