package module

import (
	"github.com/wjshen/gophrame/core/starter"

	_ "github.com/wjshen/gophrame/default/controller"
	_ "github.com/wjshen/gophrame/default/security"
	_ "github.com/wjshen/gophrame/default/service"
)

const (
	MODULE_ID = 1
)

func init() {
	starter.RegisterStarter(Start)

	// 1. 加载Config...
}

func Start() {

}
