package command

import (
	"github.com/spf13/pflag"

	"github.com/gophab/gophrame/core/global"
)

var Mode string = "production"
var Profile string = ""

// 0. 初始化
func init() {
	pflag.StringVar(&Mode, "mode", "production", "Run application in debug|production mode")
	pflag.StringVar(&Profile, "profile", "", "Run application with profile")
}

// 1. Command 解析
func Init() {
	// 1. 解析
	pflag.Parse()

	// 2.根据启动设置环境参数
	if Mode == "debug" {
		global.Debug = true
	}

	if Profile != "" {
		global.Profile = Profile
	}
}
