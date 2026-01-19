package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/security/token"
	TokenConfig "github.com/gophab/gophrame/core/security/token/config"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

type OAuth2Server struct {
	*server.Server
	once sync.Once

	AccessGenerate oauth2.AccessGenerate `inject:"accessGenerate"`
	TokenStore     oauth2.TokenStore     `inject:"tokenStore"`

	// AuthorizeGenerate oauth2.AuthorizeGenerate `inject:"authorizeGenerate"`
	UserHandler       IUserHandler       `inject:"userHandler"`
	MobileUserHandler IMobileUserHandler `inject:"userHandler"`
	EmailUserHandler  IEmailUserHandler  `inject:"userHandler"`
	SocialUserHandler ISocialUserHandler `inject:"userHandler"`

	clientAuthorizedHandlers      []server.ClientAuthorizedHandler
	clientScopeHandlers           []server.ClientScopeHandler
	passwordAuthorizationHandlers []server.PasswordAuthorizationHandler
	userAuthorizationHandlers     []server.UserAuthorizationHandler
}

var theServer *OAuth2Server

func (s *OAuth2Server) init() {
	s.once.Do(func() {
		s.initServer()
	})
}

func (s *OAuth2Server) manager() oauth2.Manager {
	// s.manager
	manager := manage.NewDefaultManager()

	manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		AccessTokenExp:     TokenConfig.Setting.AccessTokenExpireTime,
		RefreshTokenExp:    TokenConfig.Setting.RefreshTokenExpireTime,
		IsGenerateRefresh:  false,
		IsRemoveAccess:     true,
		IsRemoveRefreshing: false,
	})

	// 配置
	tokenConfig := &manage.Config{
		AccessTokenExp:    TokenConfig.Setting.AccessTokenExpireTime,
		RefreshTokenExp:   TokenConfig.Setting.RefreshTokenExpireTime,
		IsGenerateRefresh: true,
	}

	manager.SetAuthorizeCodeTokenCfg(tokenConfig)
	manager.SetImplicitTokenCfg(tokenConfig)
	manager.SetPasswordTokenCfg(tokenConfig)
	manager.SetClientTokenCfg(tokenConfig)

	// client存储方式 <= DB
	manager.MapClientStorage(ClientStore())

	// token存储方式
	// manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapTokenStorage(token.TokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapAccessGenerate(token.AccessGenerate())

	return manager
}

// 初始化服务
func (s *OAuth2Server) initServer() *server.Server {
	// oauth2Server = server.NewServer(server.NewConfig(), manager)
	s.Server = server.NewDefaultServer(s.manager())

	s.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.ClientCredentials, oauth2.PasswordCredentials, oauth2.Implicit, oauth2.Refreshing)
	s.SetAllowGetAccessRequest(true)

	// 对scope的授权，过滤非授权scope
	s.SetAuthorizeScopeHandler(s.AuthorizeScopeHandler)

	s.SetClientInfoHandler(s.ClientInfoHandler)
	s.SetClientScopeHandler(s.ClientScopeHandler)

	s.SetClientAuthorizedHandler(s.ClientAuthorizedHandler)

	// 密码授权模式才需要用到这个配置, 这个模式不需要分配授权码,而是直接分配token,通常用于无后端的应用
	s.SetPasswordAuthorizationHandler(s.PasswordAuthorizationHandler)

	// 这一行很关键,这个方法让oauth框架识别当前用户身份标识(并且可以人为处理登陆状态检验等等)
	// 具体看userAuthorizeHandler方法实现
	s.SetUserAuthorizationHandler(s.UserAuthorizationHandler)

	s.SetInternalErrorHandler(s.internalErrorHandler)
	s.SetResponseErrorHandler(s.responseErrorHandler)

	s.SetExtensionFieldsHandler(s.ExtensionFieldsHandler)

	// Test User
	s.RegisterPasswordAuthorizationHandler(s.TestPasswordAuthorizationHandler)

	// Session User
	s.RegisterUserAuthorizationHandler(s.SessionUserAuthorizationHandler)

	return s.Server
}

func (s *OAuth2Server) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Server.HandleAuthorizeRequest(w, r.WithContext(
		context.WithValue(
			context.WithValue(
				r.Context(),
				AppIdContextKey,
				r.Header.Get("X-App-Id"),
			),
			AuthorizationCodeKey,
			r.Header.Get("X-Authorization-Code"),
		)))
}

func (s *OAuth2Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Server.HandleTokenRequest(w, r.WithContext(
		context.WithValue(
			context.WithValue(
				r.Context(),
				AppIdContextKey,
				r.Header.Get("X-App-Id"),
			),
			AuthorizationCodeKey,
			r.Header.Get("X-Authorization-Code"),
		)))
}

func (s *OAuth2Server) RegisterClientScopeHandler(handler server.ClientScopeHandler) {
	if s.clientScopeHandlers == nil {
		s.clientScopeHandlers = []server.ClientScopeHandler{}
	}

	s.clientScopeHandlers = append(s.clientScopeHandlers, handler)
}

func (s *OAuth2Server) RegisterClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	if s.clientAuthorizedHandlers == nil {
		s.clientAuthorizedHandlers = []server.ClientAuthorizedHandler{}
	}

	s.clientAuthorizedHandlers = append(s.clientAuthorizedHandlers, handler)
}

func (s *OAuth2Server) RegisterPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler) {
	if s.passwordAuthorizationHandlers == nil {
		s.passwordAuthorizationHandlers = []server.PasswordAuthorizationHandler{}
	}

	s.passwordAuthorizationHandlers = append(s.passwordAuthorizationHandlers, handler)
}

func (s *OAuth2Server) RegisterUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	if s.userAuthorizationHandlers == nil {
		s.userAuthorizationHandlers = []server.UserAuthorizationHandler{}
	}

	s.userAuthorizationHandlers = append(s.userAuthorizationHandlers, handler)
}

func (s *OAuth2Server) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	if len(s.clientScopeHandlers) > 0 {
		for _, handler := range s.clientScopeHandlers {
			if allowed, err := handler(tgr); err != nil {
				return false, err
			} else if allowed {
				return allowed, nil
			}
		}
		return false, nil
	}
	return true, nil
}

func (s *OAuth2Server) ClientAuthorizedHandler(clientID string, grant oauth2.GrantType) (bool, error) {
	if len(s.clientAuthorizedHandlers) > 0 {
		for _, handler := range s.clientAuthorizedHandlers {
			if allowed, err := handler(clientID, grant); err != nil {
				return false, err
			} else if allowed {
				return allowed, nil
			}
		}
		return false, nil
	}
	return true, nil
}

// oauth框架通过本方法识别用户身份信息,并且可以人为进行登录状态校验
// 本方法正常执行后,则会为客户端分配授权码(authorization_code)
func (s *OAuth2Server) UserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if len(s.userAuthorizationHandlers) > 0 {
		for _, handler := range s.userAuthorizationHandlers {
			if userId, err := handler(w, r); err != nil {
				return "", err
			} else {
				return userId, nil
			}
		}
	}
	return "", nil
}

func (s *OAuth2Server) ExtensionFieldsHandler(ti oauth2.TokenInfo) map[string]any {
	return nil
}

func (s *OAuth2Server) SessionUserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
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

func (s *OAuth2Server) TestPasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	if username == "test" && password == "test" {
		userID = "test"
		return userID, nil
	}

	return "", nil
}

func (s *OAuth2Server) PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	if username == "test" && password == "test" {
		userID = "test"
		return userID, nil
	}

	if len(s.passwordAuthorizationHandlers) > 0 {
		for _, handler := range s.passwordAuthorizationHandlers {
			if user, err := handler(ctx, clientID, username, password); err != nil {
				return "", err
			} else if user != "" {
				return user, nil
			}
		}
	}
	return "", nil
}

func (s *OAuth2Server) WriteToken(w http.ResponseWriter, info oauth2.TokenInfo) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(s.GetTokenData(info))
}

// 根据client注册的scope过滤非法scope
func (*OAuth2Server) AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
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
