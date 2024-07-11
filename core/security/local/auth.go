package local

import (
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/security/token"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
)

type TokenValidator struct {
	TokenResolver token.ITokenResolver `inject:"tokenResolver"`
}

var tokenValidator = &TokenValidator{}

func init() {
	inject.InjectValue("tokenValidator", tokenValidator)
}

func ValidationBearerToken(ctx *gin.Context) (oauth2.TokenInfo, error) {
	// 1. 获取当前Token
	tokenValue, err := SecurityUtil.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Validate
	if tokenValidator.TokenResolver != nil {
		if ti, err := tokenValidator.TokenResolver.Resolve(ctx, tokenValue); err == nil {
			if ti != nil {
				var uid = strings.Split(ti.GetUserID(), "@")
				ctx.Set("_CURRENT_USER_ID_", uid[0])
				if len(uid) > 1 {
					ctx.Set("_CURRENT_TENANT_ID_", uid[1])
				}
			}
			return ti, err
		} else {
			return nil, err
		}
	}

	return nil, nil
}
