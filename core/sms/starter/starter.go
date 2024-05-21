package starter

import (
	_ "github.com/wjshen/gophrame/config"

	_ "github.com/wjshen/gophrame/core/sms/aliyun/starter"
	_ "github.com/wjshen/gophrame/core/sms/code/starter"
	_ "github.com/wjshen/gophrame/core/sms/qcloud/starter"
)

func init() {
}
