package mysql

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/gophab/gophrame/core/database/mysql/config"
	"github.com/gophab/gophrame/core/logger"
)

func defaultDSN() string {
	Host := config.Setting.Default.Host
	Database := config.Setting.Default.Database
	Port := config.Setting.Default.Port
	User := config.Setting.Default.User
	Password := config.Setting.Default.Password
	Charset := config.Setting.Default.Charset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", User, Password, Host, Port, Database, Charset)

	return dsn
}

func readDSN() string {
	Host := config.Setting.Read.Host
	Database := config.Setting.Read.Database
	Port := config.Setting.Read.Port
	User := config.Setting.Read.User
	Password := config.Setting.Read.Password
	Charset := config.Setting.Read.Charset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", User, Password, Host, Port, Database, Charset)

	return dsn
}

func InitDB(opts ...gorm.Option) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	logger.Info("Initialize MySQL database: ", config.Setting.Default.Host, config.Setting.Default.Database)
	db, err = openDB(defaultDSN(), opts...)
	if err != nil {
		return nil, err
	}

	if config.Setting.EnableRead {
		if readDialector := mysql.Open(readDSN()); readDialector != nil {
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
			return nil, errors.New("Open Read Dialector Error: " + readDSN())
		}
	}

	return db, nil
}

func openDB(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	dialector := mysql.Open(dsn)
	if dialector != nil {
		if db, err := gorm.Open(dialector, opts...); err != nil {
			return nil, err
		} else {
			return db, nil
		}
	}
	return nil, errors.New("Open Dialector Error: " + dsn)
}
