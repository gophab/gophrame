package sms

import (
	"sync"

	"github.com/wjshen/gophrame/core/sms/aliyun"
	"github.com/wjshen/gophrame/core/sms/qcloud"
	"github.com/wjshen/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	once.Do(func() {
		aliyun.Start()
		qcloud.Start()
	})
}
