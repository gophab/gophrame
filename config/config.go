package config

import (
	c "github.com/wjshen/gophrame/core/config"
)

var ConfigYml c.IYmlConfig = c.ConfigYml

func Init() {
	c.InitConfig()
}
