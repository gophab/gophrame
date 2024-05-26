package database

import (
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/database/config"
	MySQL "github.com/gophab/gophrame/core/database/mysql"
	"github.com/gophab/gophrame/core/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	db    *gorm.DB
	mutex sync.Mutex
)

func InitDB() *gorm.DB {
	mutex.Lock()
	if db == nil {
		logger.Info("Initializing Database Driver: ", config.Setting.Driver)
		//logger.Debug("Database configuration: ", json.String(config.Database))

		options := &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 defaultLogger(), //拦截、接管 gorm v2 自带日志
			NamingStrategy:         defaultNamingStrategy(),
		}

		switch config.Setting.Driver {
		case "mysql":
			db, _ = MySQL.InitDB(options)
		}

		// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
		// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
		_ = db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)

		// https://github.com/go-gorm/gorm/issues/4838
		_ = db.Callback().Create().Before("gorm:create").Register("UpdateCreatedTimeHook", UpdateCreatedTimeHook)
		_ = db.Callback().Create().Before("gorm:create").Register("UpdateIdHook", UpdateIdHook)
		_ = db.Callback().Update().Before("gorm:update").Register("UpdateLastModifiedTimeHook", UpdateLastModifiedTimeHook)
		_ = db.Callback().Delete().Before("gorm:delete").Register("UpdateDeletedTimeHook", UpdateDeletedTimeHook)

		// 为主连接设置连接池(43行返回的数据库驱动指针)
		if rawDb, err := db.DB(); err == nil {
			rawDb.SetConnMaxIdleTime(config.Setting.ConnectionMaxIdleTime)
			rawDb.SetConnMaxLifetime(config.Setting.ConnectionMaxLifeTime)
			rawDb.SetMaxIdleConns(config.Setting.MaxIdleConnections)
			rawDb.SetMaxOpenConns(config.Setting.MaxOpenConnections)
		}
	}
	mutex.Unlock()

	return db
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
