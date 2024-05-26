package security

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security/config"
	"github.com/gophab/gophrame/core/security/local"
	"github.com/gophab/gophrame/core/security/remote"
	"github.com/gophab/gophrame/core/security/server"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
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
