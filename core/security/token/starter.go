package token

import (
	_ "github.com/wjshen/gophrame/config"
)

func init() {
	InitTokenResolver()
	InitTokenStore()
}
