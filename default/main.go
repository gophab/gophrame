package module

import (
	"github.com/gophab/gophrame/core/starter"

	_ "github.com/gophab/gophrame/default/controller"
	_ "github.com/gophab/gophrame/default/security"
	_ "github.com/gophab/gophrame/default/service"
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
