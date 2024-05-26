package command

import (
	"github.com/spf13/pflag"

	"github.com/gophab/gophrame/core/global"
)

var Mode string = "production"
var Profile string = ""

func init() {
	pflag.StringVar(&Mode, "mode", "production", "Run application in debug|production mode")
	pflag.StringVar(&Profile, "profile", "", "Run application with profile")
	pflag.Parse()

	// 2.根据启动设置环境参数
	if Mode == "debug" {
		global.Debug = true
	}
}
