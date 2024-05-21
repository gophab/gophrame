package starter

import (
	"sync"

	_ "github.com/wjshen/gophrame/config"

	_ "github.com/wjshen/gophrame/core/email/code/starter"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {

	})
}
