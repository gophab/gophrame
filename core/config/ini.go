package config

import (
	"log"

	"github.com/go-ini/ini"
)

var iniConf *ini.File

func Init() {
	var err error
	iniConf, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Configuration.Setup, fail to parse 'conf/app.ini': %v", err)
	}
}

func MapTo(section string, v any) {
	err := iniConf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("IniConf.MapTo Setting err: %v", err)
	}
}
