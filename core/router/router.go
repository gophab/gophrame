package router

import (
	"net/http"

	_ "github.com/wjshen/gophrame/core/swagger"

	EmailCode "github.com/wjshen/gophrame/core/email/code"
	SmsCode "github.com/wjshen/gophrame/core/sms/code"

	"github.com/wjshen/gophrame/core/captcha"
	"github.com/wjshen/gophrame/core/engine"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/security/server"

	"github.com/wjshen/gophrame/core/social/dingtalk"
	"github.com/wjshen/gophrame/core/social/feishu"
	"github.com/wjshen/gophrame/core/social/wxcp"
	"github.com/wjshen/gophrame/core/social/wxmp"

	"github.com/gin-gonic/gin"
)

type RouterBootstrap struct {
	SmsCodeController   *SmsCode.SmsCodeController     `inject:"smsCodeController"`
	EmailCodeController *EmailCode.EmailCodeController `inject:"emailCodeController"`
	WxcpController      *wxcp.WxcpController           `inject:"wxcpController"`
	WxmpController      *wxmp.WxmpController           `inject:"wxmpController"`
	FeishuController    *feishu.FeishuController       `inject:"feishuController"`
	DingtalkController  *dingtalk.DingtalkController   `inject:"dingtalkController"`
}

var bootstrap = &RouterBootstrap{}

func init() {
	inject.InjectValue("routerBootstrap", bootstrap)
}

func Root() *gin.Engine {
	return engine.Get()
}

func InitDefaultRouters() *gin.Engine {
	// 处理
	Root().GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "OK")
	})

	//处理静态资源
	Root().Static("/public", "./public") //  定义静态资源路由与实际目录映射关系

	// 创建图形验证码路由
	captcha.InitRouter(Root())

	// 创建手机验证码路由
	if bootstrap.SmsCodeController != nil {
		bootstrap.SmsCodeController.InitRouter(Root())
	}

	// 创建手机验证码路由
	if bootstrap.EmailCodeController != nil {
		bootstrap.EmailCodeController.InitRouter(Root())
	}

	// Social平台
	// 企业微信
	if bootstrap.WxcpController != nil {
		bootstrap.WxcpController.InitRouter(Root())
	}

	// 微信公众号
	if bootstrap.WxmpController != nil {
		bootstrap.WxmpController.InitRouter(Root())
	}

	// 飞书
	if bootstrap.FeishuController != nil {
		bootstrap.FeishuController.InitRouter(Root())
	}

	// 钉钉
	if bootstrap.DingtalkController != nil {
		bootstrap.DingtalkController.InitRouter(Root())
	}

	// 初始化OAuth2 API路由
	server.InitRouter(Root())

	return Root()
}
