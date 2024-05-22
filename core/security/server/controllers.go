package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wjshen/gophrame/errors"

	"github.com/wjshen/gophrame/core/captcha"
	"github.com/wjshen/gophrame/core/eventbus"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/redis"
	"github.com/wjshen/gophrame/core/security/model"
	"github.com/wjshen/gophrame/core/security/server/config"
	"github.com/wjshen/gophrame/core/security/token"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-session/session"
	"github.com/patrickmn/go-cache"
)

type LoginForm struct {
	Mode     string `form:"mode" json:"mode"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 前端接口
		g.POST("/openapi/login", captcha.HandleCaptchaVerify(false), oauth2Controller.Login) // 登录

		// 后端接口
		g.GET("/openapi/auth", oauth2Controller.Auth) // 授权页面,选择需要授权的权限项

		// 增加OAuth2 Server API
		g.POST("/openapi/authorize", oauth2Controller.Authorize)      // 获取授权码 或 implicit方式请求token
		g.POST("/openapi/token", oauth2Controller.HandleTokenRequest) // 应用程序通过此请求获取token
		g.GET("/openapi/token", oauth2Controller.QueryToken)          // 根据授权码获取token
	}
}

type OAuth2Controller struct {
	OAuth2Server *OAuth2Server     `inject:"oauth2Server"`
	TokenStore   oauth2.TokenStore `inject:"tokenStore"`
	reqCache     *cache.Cache
}

var oauth2Controller *OAuth2Controller

func init() {
	if config.Setting.Enabled {
		oauth2Controller = &OAuth2Controller{reqCache: cache.New(time.Minute*5, time.Minute*5)}
		inject.InjectValue("oauth2Controller", oauth2Controller)
	}
}

/**
 * POST /login
 */
func (o *OAuth2Controller) Login(c *gin.Context) {
	clientID, clientSecret, err := o.OAuth2Server.ClientInfoHandler(c.Request)
	if err != nil {
		response.FailMessage(c, http.StatusInternalServerError, "未知应用")
		return
	}

	var loginForm = LoginForm{Mode: "password"}
	if err := c.ShouldBind(&loginForm); err == nil {
		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			response.FailMessage(c, http.StatusInternalServerError, "会话错误")
			return
		}

		var userDetails *model.UserDetails
		switch loginForm.Mode {
		case "password": // 使用用户名/密码登录
			if c.Param("captcha") == "true" && o.OAuth2Server.UserHandler != nil {
				userDetails, err = o.OAuth2Server.UserHandler.GetUserDetails(c.Request.Context(), loginForm.Username, loginForm.Password)
			}
		case "mobile": // 使用手机/验证码登录
			if c.Param("captcha") == "true" && o.OAuth2Server.MobileUserHandler != nil {
				userDetails, err = o.OAuth2Server.MobileUserHandler.GetMobileUserDetails(c.Request.Context(), loginForm.Username, loginForm.Password)
			}
		case "email": // 使用邮箱登录
			if c.Param("captcha") == "true" && o.OAuth2Server.EmailUserHandler != nil {
				userDetails, err = o.OAuth2Server.EmailUserHandler.GetEmailUserDetails(c.Request.Context(), loginForm.Username, loginForm.Password)
			}
		case "social": // 使用社交账号登录
			if o.OAuth2Server.SocialUserHandler != nil {
				var appId = c.Request.Header.Get("X-App-Id")
				userDetails, err = o.OAuth2Server.SocialUserHandler.GetSocialUserDetails(
					context.WithValue(c.Request.Context(), AppIdContextKey, appId),
					loginForm.Username,
					loginForm.Password)
			}
		}

		if userDetails == nil || err != nil {
			response.FailMessage(c, http.StatusUnauthorized, "账号密码错误")
			return
		}

		var info oauth2.TokenInfo
		info, err = o.OAuth2Server.manager.GenerateAccessToken(c.Request.Context(), oauth2.PasswordCredentials, &oauth2.TokenGenerateRequest{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			UserID:       userDetails.UserId,
			RedirectURI:  "",
			Scope:        "app",
		})
		if info == nil || err != nil {
			response.FailMessage(c, http.StatusInternalServerError, "应用未授权")
			return
		}

		code := c.Request.Header.Get("X-Authorization-Code")
		if code != "" {
			// 绑定authorization_code
			info.SetCode(code)
			info.SetCodeCreateAt(time.Now())
			info.SetCodeExpiresIn(time.Minute * 5)
			o.OAuth2Server.TokenStore.Create(c.Request.Context(), info)
		}

		// Session
		store.Set("LoggedInUserId", info.GetUserID())
		store.Save()

		// 回写Token
		c.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(o.OAuth2Server.GetTokenData(info))

		// 发送用户登录事件
		eventbus.PublishEvent("USER_LOGIN", userDetails.UserId, map[string]string{"IP": strings.Split(c.Request.RemoteAddr, ":")[0]})
	} else {
		http.Error(c.Writer, err.Error(), http.StatusNonAuthoritativeInfo)
		return
	}

}

/**
 * GET /auth
 *
 * OK => auth.html 授权页面
 */
func (o *OAuth2Controller) Auth(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 未授权，转login页面
	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.Writer.Header().Set("Location", "/login")
		c.Writer.WriteHeader(http.StatusFound)
		return
	}

	// 已授权，加载auth.html静态页面
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"title": "auth",
	})
}

/**
 * POST /authorize
 *
 * 授权接口：
 */
func (o *OAuth2Controller) Authorize(c *gin.Context) {
	store, err := session.Start(c, c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	c.Request.Form = form
	store.Delete("ReturnUri")
	store.Save()

	err = o.OAuth2Server.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

}

/**
 * POST /token
 *
 * 处理token请求
 */
func (o *OAuth2Controller) HandleTokenRequest(c *gin.Context) {
	err := o.OAuth2Server.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Abort()
}

/**
 * GET /token
 *
 * 根据授权码获取已授权的Token
 */
func (o *OAuth2Controller) QueryToken(c *gin.Context) {
	authorityCode := request.Param(c, "code").DefaultString("")

	if authorityCode == "" {
		authorityCode = c.Request.Header.Get("X-Authorization-Code")
		if authorityCode == "" {
			response.FailCode(c, errors.INVALID_PARAMS)
			return
		}
	}

	var lock = "code:" + authorityCode
	if o.reqCache.Add(lock, "1", time.Minute*5) == nil {
		defer o.reqCache.Delete(lock)
		for times := 5; times > 0; times-- {
			ti, err := o.TokenStore.GetByCode(c.Request.Context(), authorityCode)
			if err == nil && ti != nil {
				response.Success(c, o.OAuth2Server.GetTokenData(ti))
				return
			}
			time.Sleep(time.Second)
		}
	}
	response.Success(c, "")
}

func (o *OAuth2Controller) GetTokenRedis(method, code string) (oauth2.TokenInfo, error) {
	// 从Redis内获取
	client := redis.GetOneRedisClient()

	var redisKey = "goes:token:" + method + ":" + code
	value, err := client.Execute("SETNX", redisKey+":lock")
	if err == nil {
		// 没有其他相同请求
		defer client.Execute("DEL", redisKey+":lock")
		for times := 5; times > 0; times-- {
			value, err = client.Execute("GET", redisKey)
			if err != nil && value != "" {
				break
			}
			time.Sleep(time.Second)
		}
	}

	return token.ParseToken(value.(string))
}
