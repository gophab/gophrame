package mysql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gophab/gophrame/core/database"
	"github.com/gophab/gophrame/core/database/postgres/config"
)

type PostgresDriver struct {
}

func (*PostgresDriver) DSN() string {
	Host := config.Setting.Host
	Database := config.Setting.Database
	Port := config.Setting.Port
	User := config.Setting.User
	Password := config.Setting.Password
	TimeZone := config.Setting.TimeZone

	if TimeZone == "" {
		TimeZone = "UTC"
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s", Host, Port, User, Database, Password, TimeZone)

	return dsn
}

func (*PostgresDriver) ReadDSN() string {
	Host := config.Setting.Read.Host
	Database := config.Setting.Read.Database
	Port := config.Setting.Read.Port
	User := config.Setting.Read.User
	Password := config.Setting.Read.Password
	TimeZone := config.Setting.TimeZone

	if TimeZone == "" {
		TimeZone = "UTC"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s", Host, Port, User, Database, Password, TimeZone)

	return dsn
}

func (*PostgresDriver) GetDialetor(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func init() {
	database.RegisterDriver("postgres", &PostgresDriver{})
}
