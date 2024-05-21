package config

import (
	"time"

	"github.com/wjshen/gophrame/core/code/config"
)

var Setting = &config.CodeStoreSetting{
	RequestInterval: time.Minute,
	ExpireIn:        time.Minute * 5,
	Redis:           nil,
}
