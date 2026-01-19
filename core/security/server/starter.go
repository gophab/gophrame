package server

import (
	"sync"
	"time"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security/server/config"
	"github.com/patrickmn/go-cache"
)

func Init() {
	logger.Debug("Starting OAuth2 Server: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		logger.Info("Initializing OAuth2 Server")

		theServer = &OAuth2Server{
			once:                          sync.Once{},
			clientAuthorizedHandlers:      []server.ClientAuthorizedHandler{},
			clientScopeHandlers:           []server.ClientScopeHandler{},
			passwordAuthorizationHandlers: []server.PasswordAuthorizationHandler{},
			userAuthorizationHandlers:     []server.UserAuthorizationHandler{},
		}
		inject.InjectValue("oauth2Server", theServer)

		theServer.init()

		oauth2Controller := &OAuth2Controller{reqCache: cache.New(time.Minute*5, time.Minute*5)}
		inject.InjectValue("oauth2Controller", oauth2Controller)

		controller.AddController(oauth2Controller)
	}
}
