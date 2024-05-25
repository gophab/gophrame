package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	SecurityUtil "github.com/wjshen/gophrame/core/security/util"
	"github.com/wjshen/gophrame/global"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/webservice/response"
)

var __ = struct {
	Enforcer *casbin.Enforcer `inject:"enforcer"`
}{}

func init() {
	inject.InjectValue("casbin", __)
}

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if __.Enforcer != nil {
			user := SecurityUtil.GetCurrentUser(c)
			if user == nil {
				response.Unauthorized(c, "用户未登录")
				return
			}

			//获取请求的PATH
			path := c.Request.URL.Path

			// 资源
			resource := strings.TrimPrefix(path, global.StringVar("CONTEXT_ROOT"))

			// 获取请求方法
			action := c.Request.Method

			// 获取用户的角色: PUBLIC,USER,ADMIN,OPERATOR,SUPERVISOR
			roles := user.Roles
			if len(roles) == 0 {
				// 缺省角色
				roles = []string{"PUBLIC"}
			}

			success, _ := __.Enforcer.Enforce(roles, resource, action)
			if !success {
				response.FailMessage(c, http.StatusForbidden, "权限不足")
				return
			}
		}
		c.Next()
	}
}
