package security

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"

	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/config"
	"github.com/wjshen/gophrame/core/security/local"
	"github.com/wjshen/gophrame/core/security/remote"
	"github.com/wjshen/gophrame/core/security/server"
	SecurityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/core/webservice/response"
)

const (
	ErrorsParseTokenFail   string = "解析token失败"
	ErrorsTokenBaseInfo    string = "token最基本的格式错误,请提供一个有效的token!"
	ErrorsNoAuthorization  string = "token鉴权未通过，请通过token授权接口重新获取token,"
	ErrorsRefreshTokenFail string = "token不符合刷新条件,请通过登陆接口重新获取token!"

	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "参数校验失败"
)

type (
	// ErrorHandleFunc error handling function
	ErrorHandleFunc func(*gin.Context, error)
	// Config defines the config for Session middleware
	HandlerConfig struct {
		// error handling when starting the session
		ErrorHandleFunc ErrorHandleFunc
		// keys stored in the context
		TokenKey string
		// defines a function to skip middleware.Returning true skips processing
		// the middleware.
		Skipper func(*gin.Context) bool
	}
)

var (
	// DefaultConfig is the default middleware config.
	DefaultConfig = HandlerConfig{
		ErrorHandleFunc: func(ctx *gin.Context, err error) {
			//ctx.AbortWithError(500, err)
			if err != nil {
				logger.Error("Internal Server Error: ", err.Error())
			} else {
				logger.Error("Internal Server Error: ", "unknown")
			}
		},
		TokenKey: "token",
		Skipper: func(_ *gin.Context) bool {
			return false
		},
	}
)

// CheckTokenVerify Verify the access token of the middleware
// 如果有用户信息，则设置用户信息，没有则设置为空，允许公开访问
func CheckTokenVerify(conf ...HandlerConfig) gin.HandlerFunc {
	cfg := DefaultConfig
	if len(conf) > 0 {
		cfg = conf[0]
	}

	if cfg.ErrorHandleFunc == nil {
		cfg.ErrorHandleFunc = DefaultConfig.ErrorHandleFunc
	}

	tokenKey := cfg.TokenKey
	if tokenKey == "" {
		tokenKey = DefaultConfig.TokenKey
	}

	return func(context *gin.Context) {
		if cfg.Skipper != nil && cfg.Skipper(context) {
			context.Next()
			return
		}

		context.Set("_CURRENT_USER_ID_", "")
		context.Set("_CURRENT_USER_", nil)

		// 1. 从context获取token
		token, err := SecurityUtil.GetToken(context)
		if err != nil || token == "" {
			context.Next()
			return
		}

		// 2. 判断是server校验还是local校验还是remote校验
		var tokenInfo oauth2.TokenInfo
		switch config.Setting.AuthMode {
		case "server":
			tokenInfo, err = server.ValidationBearerToken(context)
		case "local":
			tokenInfo, err = local.ValidationBearerToken(context)
		case "remote":
			tokenInfo, err = remote.ValidationBearerToken(context)
		}

		if err != nil || tokenInfo == nil {
			context.Next()
			return
		}

		context.Set("_CURRENT_USER_ID_", tokenInfo.GetUserID())

		context.Set(tokenKey, tokenInfo)

		context.Next()
	}
}

// HandleTokenVerify Verify the access token of the middleware
// 只有登录获取授权后才可以访问
func HandleTokenVerify(conf ...HandlerConfig) gin.HandlerFunc {
	cfg := DefaultConfig
	if len(conf) > 0 {
		cfg = conf[0]
	}

	if cfg.ErrorHandleFunc == nil {
		cfg.ErrorHandleFunc = DefaultConfig.ErrorHandleFunc
	}

	tokenKey := cfg.TokenKey
	if tokenKey == "" {
		tokenKey = DefaultConfig.TokenKey
	}

	return func(context *gin.Context) {
		if cfg.Skipper != nil && cfg.Skipper(context) {
			context.Next()
			return
		}

		context.Set("_CURRENT_USER_ID_", "")
		context.Set("_CURRENT_USER_", nil)

		// 1. 从context获取token
		token, err := SecurityUtil.GetToken(context)
		if err != nil || token == "" {
			cfg.ErrorHandleFunc(context, err)
			TokenErrorParam(context)
			return
		}

		// 2. 判断是server校验还是local校验还是remote校验
		var tokenInfo oauth2.TokenInfo
		switch config.Setting.AuthMode {
		case "server":
			tokenInfo, err = server.ValidationBearerToken(context)
		case "local":
			tokenInfo, err = local.ValidationBearerToken(context)
		case "remote":
			tokenInfo, err = remote.ValidationBearerToken(context)
		}

		if err != nil || tokenInfo == nil {
			cfg.ErrorHandleFunc(context, err)
			ErrorTokenAuthFail(context)
			return
		}

		context.Set("_CURRENT_USER_ID_", tokenInfo.GetUserID())

		context.Set(tokenKey, tokenInfo)

		context.Next()
	}
}

// RefreshTokenConditionCheck 刷新token条件检查中间件，针对已经过期的token，要求是token格式以及携带的信息满足配置参数即可
func RefreshTokenConditionCheck() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := SecurityUtil.GetToken(context)
		if err != nil {
			TokenErrorParam(context)
			return
		}

		if token != "" {
			//
		}

	}
}

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
