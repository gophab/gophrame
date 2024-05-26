package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/engine"
	"github.com/gophab/gophrame/core/global"
)

func UseCors() {
	engine.Get().Use(Cors())
}

// Cors 直接放行所有跨域请求并放行所有 OPTIONS 方法
func Cors() gin.HandlerFunc {
	return handleCors
}

// CorsByRules 按照配置处理跨域请求
func CorsByRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.Cors.Mode == "allow-all" {
			// 放行全部
			handleCors(c)
			return
		} else {
			// 按WhiteList设定
			handleCorsWithRule(c)
			return
		}
	}
}

func handleCors(c *gin.Context) {
	method := c.Request.Method

	// 这是允许访问所有域
	c.Header("Access-Control-Allow-Origin", "*")

	//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH, OPTIONS, UPDATE")

	//  header的类型
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Authorization, Content-Length, X-CSRF-Token, Token, session, X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, Captcha, X-Verification-Code, X-Authorization-Code, X-App-Id")

	//  允许跨域设置, 可以返回其他子段
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified, Pragma, FooBar, X-Total-Count") // 跨域关键设置 让浏览器可以解析

	//  跨域请求是否需要带cookie信息 默认设置为true
	c.Header("Access-Control-Allow-Credentials", "true")

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusAccepted)
	}
	// 处理请求
	c.Next()
}

func handleCorsWithRule(c *gin.Context) {
	whitelist := checkCors(c.GetHeader("origin"))

	// 通过检查, 添加请求头
	if whitelist != nil {
		c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
		c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
		c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
		c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
		if whitelist.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
	}

	// 严格白名单模式且未通过检查，直接拒绝处理请求
	if whitelist == nil && global.Cors.Mode == "strict-whitelist" && !(c.Request.Method == "GET" && c.Request.URL.Path == "/health") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	} else {
		// 非严格白名单模式，无论是否通过检查均放行所有 OPTIONS 方法
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}

	// 处理请求
	c.Next()
}

func checkCors(currentOrigin string) *global.CORSWhitelist {
	for _, whitelist := range global.Cors.Whitelist {
		// 遍历配置中的跨域头，寻找匹配项
		if currentOrigin == whitelist.AllowOrigin {
			return &whitelist
		}
	}
	return nil
}
