package starter

import (
	"sync"

	"github.com/wjshen/gophrame/core/social/dingtalk"
	"github.com/wjshen/gophrame/core/social/feishu"
	"github.com/wjshen/gophrame/core/social/wxcp"
	"github.com/wjshen/gophrame/core/social/wxma"
	"github.com/wjshen/gophrame/core/social/wxmp"

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
		dingtalk.Start()
		feishu.Start()
		wxcp.Start()
		wxmp.Start()
		wxma.Start()
	})
}
