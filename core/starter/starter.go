package starter

import (
	_ "github.com/wjshen/gophrame/core/casbin"
	_ "github.com/wjshen/gophrame/core/database"
	_ "github.com/wjshen/gophrame/core/destroy" // 监听程序退出信号，用于资源的释放
	_ "github.com/wjshen/gophrame/core/email/starter"
	_ "github.com/wjshen/gophrame/core/engine"
	_ "github.com/wjshen/gophrame/core/eventbus"
	_ "github.com/wjshen/gophrame/core/rabbitmq"
	_ "github.com/wjshen/gophrame/core/redis"
	_ "github.com/wjshen/gophrame/core/security"
	_ "github.com/wjshen/gophrame/core/sms/starter"
	_ "github.com/wjshen/gophrame/core/snowflake"
	_ "github.com/wjshen/gophrame/core/social/starter"
	_ "github.com/wjshen/gophrame/core/websocket"
)

func init() {

}
