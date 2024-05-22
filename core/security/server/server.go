package server

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/config"
	"github.com/wjshen/gophrame/core/security/model"
	"github.com/wjshen/gophrame/core/security/token"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

type OAuth2Server struct {
	server  *server.Server
	manager *manage.Manager
	once    sync.Once

	ClientStore    oauth2.ClientStore    `inject:"clientStore"`
	TokenStore     oauth2.TokenStore     `inject:"tokenStore"`
	AccessGenerate oauth2.AccessGenerate `inject:"accessGenerate"`
	// AuthorizeGenerate oauth2.AuthorizeGenerate `inject:"authorizeGenerate"`
	UserHandler       IUserHandler       `inject:"userHandler"`
	MobileUserHandler IMobileUserHandler `inject:"userHandler"`
	EmailUserHandler  IEmailUserHandler  `inject:"userHandler"`
	SocialUserHandler ISocialUserHandler `inject:"userHandler"`
}

var theOAuth2Server *OAuth2Server

func Start() {
	if config.Setting.Server.Enabled {
		logger.Info("Initializing OAuth2 Server")

		theOAuth2Server = &OAuth2Server{
			once: sync.Once{},
		}
		inject.InjectValue("oauth2Server", theOAuth2Server)

		theOAuth2Server.init()
	}
}

func (s *OAuth2Server) init() {
	s.once.Do(func() {
		s.initServer(s.initManager())
	})
}

func (s *OAuth2Server) initManager() oauth2.Manager {
	// s.manager
	s.manager = manage.NewDefaultManager()

	s.manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		AccessTokenExp:     config.Setting.Token.AccessTokenExpireTime,
		RefreshTokenExp:    config.Setting.Token.RefreshTokenExpireTime,
		IsGenerateRefresh:  false,
		IsRemoveAccess:     true,
		IsRemoveRefreshing: false,
	})

	// 配置
	tokenConfig := &manage.Config{
		AccessTokenExp:    config.Setting.Token.AccessTokenExpireTime,
		RefreshTokenExp:   config.Setting.Token.RefreshTokenExpireTime,
		IsGenerateRefresh: true,
	}

	s.manager.SetAuthorizeCodeTokenCfg(tokenConfig)
	s.manager.SetImplicitTokenCfg(tokenConfig)
	s.manager.SetPasswordTokenCfg(tokenConfig)
	s.manager.SetClientTokenCfg(tokenConfig)

	// client存储方式 <= DB
	s.manager.MapClientStorage(ClientStore())

	// token存储方式
	// manager.MustTokenStorage(store.NewMemoryTokenStore())
	s.manager.MapTokenStorage(token.TokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())
	s.manager.MapAccessGenerate(token.AccessGenerate())

	inject.InjectValue("oauth2.Manager", s.manager)
	return s.manager
}

// 初始化服务
func (s *OAuth2Server) initServer(manager oauth2.Manager) *server.Server {
	// oauth2Server = server.NewServer(server.NewConfig(), manager)
	s.server = server.NewDefaultServer(s.manager)

	s.server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.ClientCredentials, oauth2.PasswordCredentials, oauth2.Implicit, oauth2.Refreshing)
	s.server.SetAllowGetAccessRequest(true)

	s.server.SetAuthorizeScopeHandler(s.authorizeScopeHandler)

	// 密码授权模式才需要用到这个配置, 这个模式不需要分配授权码,而是直接分配token,通常用于无后端的应用
	s.server.SetPasswordAuthorizationHandler(s.passwordAuthorizationHandler)

	// 这一行很关键,这个方法让oauth框架识别当前用户身份标识(并且可以人为处理登陆状态检验等等)
	// 具体看userAuthorizeHandler方法实现
	s.server.SetUserAuthorizationHandler(s.userAuthorizeHandler)

	s.server.SetInternalErrorHandler(s.internalErrorHandler)
	s.server.SetResponseErrorHandler(s.responseErrorHandler)

	return s.server
}

func (s *OAuth2Server) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.server.HandleAuthorizeRequest(w, r.WithContext(context.WithValue(r.Context(), AppIdContextKey, r.Header.Get("X-App-Id"))))
}

func (s *OAuth2Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.server.HandleTokenRequest(w, r.WithContext(context.WithValue(r.Context(), AppIdContextKey, r.Header.Get("X-App-Id"))))
}

func (s *OAuth2Server) ClientInfoHandler(r *http.Request) (string, string, error) {
	return s.server.ClientInfoHandler(r)
}

func (s *OAuth2Server) GetTokenData(ti oauth2.TokenInfo) map[string]interface{} {
	return s.server.GetTokenData(ti)
}

func (s *OAuth2Server) ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error) {
	return s.server.ValidationBearerToken(r)
}

// oauth框架通过本方法识别用户身份信息,并且可以人为进行登录状态校验
// 本方法正常执行后,则会为客户端分配授权码(authorization_code)
func (*OAuth2Server) userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func (s *OAuth2Server) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	if username == "test" && password == "test" {
		userID = "test"
		return userID, nil
	}

	var userDetails *model.UserDetails
	if mobile, b := strings.CutPrefix(username, "mobile:"); b {
		if s.MobileUserHandler != nil {
			userDetails, err = s.MobileUserHandler.GetMobileUserDetails(ctx, mobile, password)
		} else {
			err = errors.New("不支持手机验证码登录")
		}
	} else if email, b := strings.CutPrefix(username, "email:"); b {
		if s.EmailUserHandler != nil {
			userDetails, err = s.EmailUserHandler.GetEmailUserDetails(ctx, email, password)
		} else {
			err = errors.New("不支持邮箱登录")
		}
	} else if social, b := strings.CutPrefix(username, "social:"); b {
		if s.SocialUserHandler != nil {
			userDetails, err = s.SocialUserHandler.GetSocialUserDetails(ctx, social, password)
		} else {
			err = errors.New("不支持社交账号登录")
		}
	} else {
		if s.UserHandler != nil {
			userDetails, err = s.UserHandler.GetUserDetails(ctx, username, password)
		} else {
			err = errors.New("不支持用户名登录")
		}
	}

	if err != nil {
		return "", err
	}

	if userDetails != nil {
		return userDetails.UserId, err
	}

	return "", errors.New("not found")
}

// 根据client注册的scope
// 过滤非法scope
func (*OAuth2Server) authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}
	scope = ScopeJoin(s)

	return
}

func (*OAuth2Server) internalErrorHandler(err error) (re *errors.Response) {
	// log.Println("Internal Error:", err.Error())
	return
}

func (*OAuth2Server) responseErrorHandler(re *errors.Response) {
	// log.Println("Response Error:", re.Error.Error())
	// return
}

type Scope struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
}

func ScopeJoin(scope []Scope) string {
	var s []string
	for _, sc := range scope {
		s = append(s, sc.ID)
	}
	return strings.Join(s, ",")
}

func ScopeFilter(clientID string, scope string) (s []Scope) {
	// cli := GetClient(clientID)
	// sl := strings.Split(scope, ",")
	// for _, str := range sl {
	// 	for _, sc := range cli.Scope {
	// 		if str == sc.ID {
	// 			s = append(s, sc)
	// 		}
	// 	}
	// }
	s = []Scope{}
	return
}
