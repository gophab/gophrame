package webservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这是允许访问所有域
		c.Header("Access-Control-Allow-Origin", "*")

		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH, OPTIONS, UPDATE")

		//  header的类型
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Authorization, Content-Length, X-CSRF-Token, Token, session, X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, Captcha, X-Verification-Code, X-Authorization-Code, X-App-Id")

		//  允许跨域设置, 可以返回其他子段
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified, Pragma, FooBar, X-Total-Count") // 跨域关键设置 让浏览器可以解析

		//  跨域请求是否需要带cookie信息 默认设置为true
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
		}

		// 处理请求
		c.Next() //  处理请求
	}
}
