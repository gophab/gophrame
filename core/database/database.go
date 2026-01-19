package database

import (
	"errors"
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/database/config"

	"github.com/gophab/gophrame/core/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type DatabaseDriver interface {
	DSN() string
	ReadDSN() string
	GetDialetor(dsn string) gorm.Dialector
}

var (
	db      *gorm.DB
	mutex   sync.Mutex
	drivers map[string]DatabaseDriver = make(map[string]DatabaseDriver)
)

func RegisterDriver(name string, driver DatabaseDriver) {
	drivers[name] = driver
}

func openDB(driver DatabaseDriver, opts ...gorm.Option) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	logger.Info("Initialize Database: ", config.Setting.Driver)
	if dialector := driver.GetDialetor(driver.DSN()); dialector != nil {
		if db, err = gorm.Open(dialector, opts...); err != nil {
			return nil, err
		}
	}
	if db == nil {
		logger.Error("Initialize Database error: ", config.Setting.Driver)
		return nil, errors.New("ERROR")
	}

	if config.Setting.Read != nil && config.Setting.Read.Enabled {
		dsn := driver.ReadDSN()
		if dsn == "" {
			dsn = driver.DSN()
		}
		if readDialector := driver.GetDialetor(dsn); readDialector != nil {
			resolverConf := dbresolver.Config{
				Replicas: []gorm.Dialector{readDialector}, //  读 操作库，查询类
				Policy:   dbresolver.RandomPolicy{},       // sources/replicas 负载均衡策略适用于
			}

			if err = db.Use(dbresolver.Register(resolverConf).
				SetConnMaxIdleTime(config.Setting.Read.ConnectionMaxIdleTime).
				SetConnMaxLifetime(config.Setting.Read.ConnectionMaxLifeTime).
				SetMaxIdleConns(config.Setting.Read.MaxIdleConnections).
				SetMaxOpenConns(config.Setting.Read.MaxOpenConnections)); err != nil {
				return nil, err
			}
		} else {
			logger.Error("Open Read Dialector Error: ", driver.ReadDSN())
		}
	}

	return db, nil
}

func InitDB() (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	if db == nil {
		logger.Info("Initializing Database Driver: ", config.Setting.Driver)
		//logger.Debug("Database configuration: ", json.String(config.Database))

		options := &gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   defaultLogger(), //拦截、接管 gorm v2 自带日志
			NamingStrategy:                           defaultNamingStrategy(),
		}

		if driver, b := drivers[config.Setting.Driver]; b && driver != nil {
			if db, err = openDB(driver, options); err != nil {
				logger.Error("Init database error: ", config.Setting.Driver, err.Error())
				return nil, err
			}

			if db == nil {
				logger.Error("Init database error: ", config.Setting.Driver)
				return nil, errors.New("ERROR")
			}
			// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
			// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
			_ = db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)

			// https://github.com/go-gorm/gorm/issues/4838
			// _ = db.Callback().Create().Before("gorm:create").Register("UpdateCreatedTimeHook", UpdateCreatedTimeHook)
			// _ = db.Callback().Create().Before("gorm:create").Register("UpdateIdHook", UpdateIdHook)
			// _ = db.Callback().Update().Before("gorm:update").Register("UpdateLastModifiedTimeHook", UpdateLastModifiedTimeHook)
			// _ = db.Callback().Delete().Before("gorm:delete").Register("UpdateDeletedTimeHook", UpdateDeletedTimeHook)

			// 为主连接设置连接池(43行返回的数据库驱动指针)
			if rawDb, err := db.DB(); err == nil {
				rawDb.SetConnMaxIdleTime(config.Setting.ConnectionMaxIdleTime)
				rawDb.SetConnMaxLifetime(config.Setting.ConnectionMaxLifeTime)
				rawDb.SetMaxIdleConns(config.Setting.MaxIdleConnections)
				rawDb.SetMaxOpenConns(config.Setting.MaxOpenConnections)
			}
		} else {
			logger.Error("Unsupported database: ", config.Setting.Driver)
			return nil, errors.New("ERROR")
		}
	}

	return db, nil
}

type MyNamingStrategy struct {
	schema.NamingStrategy
}

func (ns MyNamingStrategy) TableName(str string) string {
	if strings.HasPrefix(str, "sys_") || strings.HasPrefix(str, "auth_") {
		return str
	}
	return ns.NamingStrategy.TableName(str)
}

func (ns MyNamingStrategy) JoinTableName(str string) string {
	if strings.HasPrefix(str, "sys_") || strings.HasPrefix(str, "auth_") {
		return str
	}
	return ns.NamingStrategy.JoinTableName(str)
}

func defaultNamingStrategy() schema.Namer {
	return &MyNamingStrategy{
		schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.Setting.TablePrefix,
		},
	}
}

func DB() *gorm.DB {
	if db == nil {
		if config.Setting.Enabled {
			InitDB()
		} else {
			logger.Error("Database Not Enabled")
		}
	}
	return db
}

func CloseDB() {
}
