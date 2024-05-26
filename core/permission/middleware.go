package permission

import (
	"net/http"
	"strings"

	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/inject"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gin-gonic/gin"
)

var __ = struct {
	PermissionService PermissionService `inject:"permissionService"`
}{}

func init() {
	inject.InjectValue("permission", __)
}

// 系统用户可以访问
func NeedSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId := SecurityUtil.GetCurrentTenantId(c)
		if tenantId == "" {
			response.Unauthorized(c, "用户未登录")
		} else if tenantId != "SYSTEM" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"data": "登录用户没有权限",
				"msg":  "ok",
			})
		} else {
			c.Next()
		}
	}
}

// 管理员可以访问
func NeedAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := SecurityUtil.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"data": "未登录用户",
				"msg":  "ok",
			})
			c.Abort()
		} else if !user.Admin {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"data": "登录用户没有权限",
				"msg":  "ok",
			})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// casbin检查用户对应的角色权限是否允许访问接口
func CheckUserRole(roles ...string) gin.HandlerFunc {
	if len(roles) <= 0 {
		return func(c *gin.Context) { c.Abort() }
	}
	var checkRoles = make(map[string]bool)
	for _, role := range roles {
		checkRoles[role] = true
	}
	return func(c *gin.Context) {
		user := SecurityUtil.GetCurrentUser(c)
		if user == nil || user.UserId == nil {
			response.Unauthorized(c, "用户未登录")
			return
		}

		for _, role := range user.Roles {
			if _, b := checkRoles[role]; b {
				c.Next()
				return
			}
		}
		response.ErrorMessage(c, http.StatusForbidden, http.StatusMethodNotAllowed, "登录用户没有权限")
	}
}

// casbin检查用户对应的角色权限是否允许访问接口
func CheckUserPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if __.PermissionService == nil {
			c.Abort()
			return
		}

		userId := SecurityUtil.GetCurrentUserId(c)
		if userId == "" {
			response.Unauthorized(c, "用户未登录")
			return
		}

		//获取请求的PATH
		// 资源
		resource := strings.TrimPrefix(c.Request.URL.Path, global.StringVar("CONTEXT_ROOT"))
		// 获取请求方法
		action := c.Request.Method
		// 从权限服务中检查接口调用权限
		ok, err := __.PermissionService.CheckPermission(userId, resource, action)
		if err != nil {
			response.Unauthorized(c, "用户未登录")
			return
		} else if !ok {
			response.ErrorMessage(c, http.StatusForbidden, http.StatusMethodNotAllowed, "登录用户没有权限")
			return
		} else {
			c.Next()
		}
	}
}
