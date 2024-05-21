package securityUtil

import (
	"fmt"
	"strings"

	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/errors"
	"github.com/wjshen/gophrame/service"

	"github.com/gin-gonic/gin"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

func GetToken(c *gin.Context) (string, error) {
	code := errors.SUCCESS

	headerParams := HeaderParams{}

	//  推荐使用 ShouldBindHeader 方式获取头参数
	if err := c.ShouldBindHeader(&headerParams); err != nil {
		code = errors.INVALID_PARAMS
	} else {
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) > 1 {
			return token[1], nil
		} else {
			code = errors.ERROR_AUTH_CHECK_TOKEN_FAIL
		}
	}

	return "", fmt.Errorf("authorization error: %d %s", code, errors.GetErrorMessage(code))
}

func GetCurrentUserId(c *gin.Context) string {
	currentUserId := c.Value("_CURRENT_USER_ID_").(string)
	if currentUserId != "" {
		return currentUserId
	}
	return ""
}

func GetCurrentTenantId(c *gin.Context) string {
	currentTenantId := c.Value("_CURRENT_TENANT_ID_").(string)
	if currentTenantId != "" {
		return currentTenantId
	}

	currentUser := GetCurrentUser(c)
	if currentUser != nil {
		c.Set("_CURRENT_TENANT_ID_", currentUser.TenantId)
		return currentUser.TenantId
	}
	return ""
}

func GetCurrentUser(c *gin.Context) *domain.User {
	if c.Value("_CURRENT_USER_") != nil {
		return c.Value("_CURRENT_USER_").(*domain.User)
	}

	currentUserId := GetCurrentUserId(c)
	if currentUserId != "" {
		if strings.HasPrefix(currentUserId, "sns:") {
			// 社交账户
			if socialUser, _ := service.GetSocialUserService().GetById(strings.SplitN(currentUserId, ":", 2)[1]); socialUser != nil {
				if socialUser.UserId == nil {
					// 未绑定系统账号
					currentUser := &domain.User{
						UserBase: domain.UserBase{
							Entity: domain.Entity{
								Id: "sns:" + socialUser.Id,
							},
							DelFlag:       socialUser.DelFlag,
							Name:          socialUser.Name,
							Mobile:        socialUser.Mobile,
							Email:         socialUser.Email,
							Avatar:        socialUser.Avatar,
							LoginTimes:    socialUser.LoginTimes,
							LastLoginTime: socialUser.LastLoginTime,
							LastLoginIp:   socialUser.LastLoginIp,
						},
					}
					c.Set("_CURRENT_USER_", currentUser)
					return currentUser
				} else {
					// 已绑定系统账号，则返回系统账户
					currentUserId = *socialUser.UserId
				}
			} else {
				return nil
			}
		}

		if currentUserId != "" {
			currentUser, _ := service.GetUserService().GetById(currentUserId)
			c.Set("_CURRENT_USER_", currentUser)
			return currentUser
		}
	}

	return nil
}
