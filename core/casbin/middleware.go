package casbin

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/response"
)

// CasbinHandler 拦截器
// Casbin检查用户对应的角色权限是否允许访问接口

var __ = struct {
	Enforcer *casbin.SyncedEnforcer `inject:"enforcer"`
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
			// 资源
			resource := strings.TrimPrefix(c.Request.URL.Path, global.StringVar("CONTEXT_ROOT"))

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
				response.ErrorMessage(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, "登录用户没有权限")
				return
			}
		}
		c.Next()
	}
}
