package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/permission"
	"github.com/gophab/gophrame/core/security"
	"github.com/gophab/gophrame/core/server/config"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/webservice/middleware"

	_ "github.com/gophab/gophrame/controller/security"
)

var ApiResources = &controller.Controllers{
	Base: "/api",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(), // oauth2 验证
	},
	Controllers: []controller.Controller{},
}

var InternalResources = &controller.Controllers{
	Base: "/api/_",
	Handlers: []gin.HandlerFunc{
		security.CheckAuthCode("_!@#$QWERasdfzxcv_"), // oauth2 验证
	},
	Controllers: []controller.Controller{},
}

var MApiResources = &controller.Controllers{
	Base: "/mapi",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),      // oauth2 验证
		permission.NeedSystemUser(),       // 需要系统用户
		permission.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{},
}

var PublicResources = &controller.Controllers{
	Base: "/openapi/public",
	Handlers: []gin.HandlerFunc{
		security.CheckTokenVerify(), // oauth2 验证
	},
	Controllers: []controller.Controller{},
}

var UserResources = &controller.Controllers{
	Base: "/openapi/user",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),      // oauth2 验证
		permission.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{},
}

var AdminResources = &controller.Controllers{
	Base: "/openapi/admin",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(), // oauth2 验证
		permission.NeedAdmin(),
		permission.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{},
}

var OpenApiResources = &controller.Controllers{
	Base: "/openapi",
	Handlers: []gin.HandlerFunc{
		security.HandleTokenVerify(),      // oauth2 验证
		permission.CheckUserPermissions(), // 权限验证
	},
	Controllers: []controller.Controller{},
}

var Resources = map[string]*controller.Controllers{
	"/api":            ApiResources,
	"/api/_":          InternalResources,
	"/mapi":           MApiResources,
	"/openapi":        OpenApiResources,
	"/openapi/public": PublicResources,
	"/openapi/user":   UserResources,
	"/openapi/admin":  AdminResources,
}

func AddController(c controller.Controller) {
	controller.AddController(c)
}

func AddControllers(cs ...controller.Controller) {
	controller.AddControllers(cs...)
}

func AddSchemaControllers(schema string, cs ...controller.Controller) {
	if resources, b := Resources[schema]; b {
		resources.AddController(cs...)
	} else {
		controllers := &controller.Controllers{
			Base:        schema,
			Controllers: append([]controller.Controller{}, cs[:]...),
		}
		Resources[schema] = controllers
		controller.AddControllers(controllers)
	}
}

func init() {
	starter.RegisterInitializor(Init)
	starter.RegisterStarter(Start)
}

// Auto Initialize entrypoint
func Init() {
	controller.AddControllers(ApiResources, InternalResources, MApiResources, OpenApiResources, PublicResources, UserResources, AdminResources)
}

// Autostart entrypoint
func Start() {
	logger.Info("Starting Framework Controllers...")
	logger.Info("Allow Cross Domain: ", config.Setting.AllowCrossDomain)
	if config.Setting.AllowCrossDomain {
		middleware.UseCors()
	}
}
