package security

import (
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/webservice/response"
)

const (
	ErrorsParseTokenFail   string = "解析token失败"
	ErrorsTokenBaseInfo    string = "token最基本的格式错误,请提供一个有效的token!"
	ErrorsNoAuthorization  string = "token鉴权未通过，请通过token授权接口重新获取token,"
	ErrorsRefreshTokenFail string = "token不符合刷新条件,请通过登陆接口重新获取token!"

	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "参数校验失败"
)

// token 基本的格式错误
func ErrorTokenBaseInfo(c *gin.Context) {
	response.Unauthorized(c, ErrorsTokenBaseInfo)
}

// token 权限校验失败
func ErrorTokenAuthFail(c *gin.Context) {
	response.Unauthorized(c, ErrorsNoAuthorization)
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// token 不符合刷新条件
func ErrorTokenRefreshFail(c *gin.Context) {
	response.Unauthorized(c, ErrorsRefreshTokenFail)
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// token 参数校验错误
func TokenErrorParam(c *gin.Context) {
	response.Unauthorized(c, ValidatorParamsCheckFailMsg)
	c.Abort()
}
