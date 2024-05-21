package starter

import (
	_ "github.com/wjshen/gophrame/config"

	_ "github.com/wjshen/gophrame/core/social/dingtalk/starter"
	_ "github.com/wjshen/gophrame/core/social/feishu/starter"
	_ "github.com/wjshen/gophrame/core/social/wxcp/starter"
	_ "github.com/wjshen/gophrame/core/social/wxma/starter"
	_ "github.com/wjshen/gophrame/core/social/wxmp/starter"
)

func init() {
}
