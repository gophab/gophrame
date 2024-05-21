package security

import (
	"net/http"

	"github.com/wjshen/gophrame/core/global"
	SecurityUtil "github.com/wjshen/gophrame/core/security/util"

	"github.com/wjshen/gophrame/service"

	"github.com/gin-gonic/gin"
)

var (
	initialized bool = false
)

// 系统用户可以访问
func NeedSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId := SecurityUtil.GetCurrentTenantId(c)
		if tenantId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"data": "未登录用户",
				"msg":  "ok",
			})
			c.Abort()
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
func CheckUserPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.Enforcer == nil {
			c.Next()
			return
		}

		userId := SecurityUtil.GetCurrentUserId(c)
		if userId == "" {
			c.Next()
			return
		}

		_ = loadCasbinPolicyData()

		requstUrl := c.Request.URL.Path
		method := c.Request.Method

		// 用户角色id需要存储在缓存，加快接口验证的效率(2021-03-11  后续实现)
		isPass, err := global.Enforcer.Enforce(userId, requstUrl, method)
		//fmt.Printf("Casbin权限校验参数打印：isPass:%v, 角色ID：%d ,url：%s ,method: %s\n", isPass, roleId, requstUrl, method)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusOK,
				"data": err,
				"msg":  "ok",
			})
			c.Abort()
			return
		} else if !isPass {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusForbidden,
				"data": "登录用户没有权限",
				"msg":  "ok",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func loadCasbinPolicyData() error {
	if !initialized {
		initialized = true

		err := service.GetRoleService().LoadAllPolicy()
		if err != nil {
			return err
		}

		err = service.GetUserService().LoadAllPolicy()
		if err != nil {
			return err
		}
	}
	return nil
}
