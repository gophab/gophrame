package SecurityUtil

import (
	"fmt"
	"strings"

	SecurityModel "github.com/gophab/gophrame/core/security/model"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/service"

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
		c.Set("_CURRENT_TENANT_ID_", util.StringValue(currentUser.TenantId))
		return util.StringValue(currentUser.TenantId)
	}
	return ""
}

func GetCurrentUser(c *gin.Context) *SecurityModel.UserDetails {
	if c.Value("_CURRENT_USER_") != nil {
		return c.Value("_CURRENT_USER_").(*SecurityModel.UserDetails)
	}

	currentUserId := GetCurrentUserId(c)
	if currentUserId != "" {
		if strings.HasPrefix(currentUserId, "sns:") {
			// 社交账户
			if socialUser, _ := service.GetSocialUserService().GetById(strings.SplitN(currentUserId, ":", 2)[1]); socialUser != nil {
				if socialUser.UserId == nil {
					// 未绑定系统账号
					userDetail := &SecurityModel.UserDetails{
						SocialId: socialUser.Id,
						Name:     socialUser.Name,
						Mobile:   socialUser.Mobile,
						Email:    socialUser.Email,
						Avatar:   socialUser.Avatar,
					}
					c.Set("_CURRENT_USER_", userDetail)
					return userDetail
				} else {
					// 已绑定系统账号，则返回系统账户
					currentUserId = *socialUser.UserId
				}
			} else {
				return nil
			}
		}

		if currentUserId != "" {
			if service.GetUserService() != nil {
				if currentUser, err := service.GetUserService().GetById(currentUserId); err == nil {
					userDetails := &SecurityModel.UserDetails{
						UserId:   currentUser.Id,
						Login:    currentUser.Login,
						Mobile:   currentUser.Mobile,
						Email:    currentUser.Email,
						SocialId: currentUser.SocialId,
						TenantId: currentUser.TenantId,
					}
					c.Set("_CURRENT_USER_", userDetails)
					return userDetails
				}
			}
			return &SecurityModel.UserDetails{
				UserId: &currentUserId,
			}
		}
	}

	return nil
}
