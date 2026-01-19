package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gophab/gophrame/core/database"
	"github.com/gophab/gophrame/core/database/mysql/config"
)

type MysqlDriver struct{}

func (*MysqlDriver) DSN() string {
	Host := config.Setting.Host
	Database := config.Setting.Database
	Port := config.Setting.Port
	User := config.Setting.User
	Password := config.Setting.Password
	Charset := config.Setting.Charset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", User, Password, Host, Port, Database, Charset)

	return dsn
}

func (*MysqlDriver) ReadDSN() string {
	if config.Setting.Read == nil {
		return ""
	}

	Host := config.Setting.Read.Host
	Database := config.Setting.Read.Database
	Port := config.Setting.Read.Port
	User := config.Setting.Read.User
	Password := config.Setting.Read.Password
	Charset := config.Setting.Read.Charset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", User, Password, Host, Port, Database, Charset)

	return dsn
}

func (*MysqlDriver) GetDialetor(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func init() {
	database.RegisterDriver("mysql", &MysqlDriver{})
}
