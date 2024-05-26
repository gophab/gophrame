package token

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security/token/config"
	JWT "github.com/gophab/gophrame/core/security/token/jwt"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	theAccessGenerate oauth2.AccessGenerate
)

func AccessGenerate() oauth2.AccessGenerate {
	if theAccessGenerate == nil {
		theAccessGenerate = InitTokenGenerator()
	}
	return theAccessGenerate
}

func init() {
	InitTokenGenerator()
}

func InitTokenGenerator() oauth2.AccessGenerate {
	if theAccessGenerate == nil {
		var result oauth2.AccessGenerate
		if config.Setting.UseJwtToken {
			result = NewJWTAccessGenerate()
		} else {
			result = NewUUIDAccessGenerate()
		}

		if result != nil {
			inject.InjectValue("accessGenerate", result)
		}

		theAccessGenerate = result
	}

	return theAccessGenerate
}

/**
 * JWT Token Generator
 */
type JWTTokenGenerator struct {
}

func NewJWTAccessGenerate() oauth2.AccessGenerate {
	logger.Debug("Using JWT token generator")
	return &JWTTokenGenerator{}
}

func (a *JWTTokenGenerator) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	claims := &JWT.Claims{
		Scope:       data.TokenInfo.GetScope(),
		RedirectURI: data.TokenInfo.GetRedirectURI(),
		StandardClaims: jwt.StandardClaims{
			Audience:  data.Client.GetID(),
			Subject:   data.UserID,
			IssuedAt:  data.TokenInfo.GetAccessCreateAt().Unix(),
			ExpiresAt: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		},
	}

	access, err := JWT.GenerateToken(claims)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

/**
 * UUID Token Generator
 */
func NewUUIDAccessGenerate() oauth2.AccessGenerate {
	logger.Debug("Using UUID token generator")
	return &UUIDTokenGenerator{}
}

type UUIDTokenGenerator struct {
	TokenStore ITokenStore `inject:"tokenStore"`
}

func (g *UUIDTokenGenerator) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	var authentication Authentication = Authentication{
		UserId:        data.UserID,
		Authenticated: true,
		Request: TokenRequest{
			ClientId: data.TokenInfo.GetClientID(),
			Scope:    data.TokenInfo.GetScope(),
		},
	}

	access := uuid.NewString()
	refresh := uuid.NewString()

	if config.Setting.ReuseAccessToken {
		// authentication => Token
		token, err := g.TokenStore.GetToken(ctx, authentication.GetId())
		if err != nil {
			return "", "", err
		}

		if token != nil && !isRefreshExpired(token) && config.Setting.ReuseRefreshToken {
			// 1. RefreshToekn is not expired, reuse refresh
			refresh = token.GetRefresh()

			if !isExpired(token) && config.Setting.ReuseAccessToken {
				// 2. AccessToken is not expired, reuse access
				access = token.GetAccess()
			}
		}
	}

	return access, refresh, nil
}

func isExpired(token oauth2.TokenInfo) bool {
	return time.Now().After(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()))
}

func isRefreshExpired(token oauth2.TokenInfo) bool {
	return time.Now().After(token.GetRefreshCreateAt().Add(token.GetRefreshExpiresIn()))
}
