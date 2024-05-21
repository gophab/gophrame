package token

import (
	"context"
	"time"

	_ "github.com/wjshen/gophrame/config"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/token/config"
	JWT "github.com/wjshen/gophrame/core/security/token/jwt"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

var (
	theTokenResolver ITokenResolver
)

func TokenResolver() ITokenResolver {
	if theTokenResolver == nil {
		theTokenResolver = InitTokenResolver()
	}
	return theTokenResolver
}

func init() {
	InitTokenResolver()
}

func InitTokenResolver() ITokenResolver {
	if theTokenResolver == nil {
		var result ITokenResolver
		if config.Setting.UseJwtToken {
			result = NewJWTTokenResolver()
		} else {
			result = NewStoreTokenResolver()
		}

		inject.InjectValue("tokenResolver", result)

		theTokenResolver = result
	}

	return theTokenResolver
}

type ITokenResolver interface {
	Resolve(context.Context, string) (oauth2.TokenInfo, error)
}

/**
 * JWT Token Resolver
 */
func NewJWTTokenResolver() ITokenResolver {
	logger.Debug("Using JWT token resolver")
	return &JWTTokenResolver{}
}

type JWTTokenResolver struct{}

//	type StandardClaims struct {
//		Audience  string `json:"aud,omitempty"`
//		ExpiresAt int64  `json:"exp,omitempty"`	AccessCreateAt + AccessExpiresIn
//		Id        string `json:"jti,omitempty"`
//		IssuedAt  int64  `json:"iat,omitempty"` AccessCreateAt
//		Issuer    string `json:"iss,omitempty"`
//		NotBefore int64  `json:"nbf,omitempty"`
//		Subject   string `json:"sub,omitempty"` UserId
//	}
func (v *JWTTokenResolver) Resolve(ctx context.Context, tokenValue string) (oauth2.TokenInfo, error) {
	if claim, _ := JWT.ParseToken(tokenValue); claim != nil {
		if err := claim.Valid(); err == nil {
			return &models.Token{
				UserID:           claim.Subject,
				ClientID:         claim.Audience,
				Code:             tokenValue,
				CodeCreateAt:     time.Unix(claim.IssuedAt, 0),
				CodeExpiresIn:    time.Second * time.Duration(claim.ExpiresAt-claim.IssuedAt),
				Access:           tokenValue,
				AccessCreateAt:   time.Unix(claim.IssuedAt, 0),
				AccessExpiresIn:  time.Second * time.Duration(claim.ExpiresAt-claim.IssuedAt),
				Refresh:          tokenValue,
				RefreshCreateAt:  time.Unix(claim.IssuedAt, 0),
				RefreshExpiresIn: time.Second * time.Duration(claim.ExpiresAt-claim.IssuedAt),
				Scope:            claim.Scope,
				RedirectURI:      claim.RedirectURI,
			}, nil
		} else {
			return nil, err
		}
	}
	return nil, nil
}

/**
 * Store Token Resolver
 */
func NewStoreTokenResolver() ITokenResolver {
	logger.Debug("Using store token resolver")
	return &StoreTokenResolver{}
}

type StoreTokenResolver struct {
	TokenStore oauth2.TokenStore `inject:"tokenStore"`
}

func (v *StoreTokenResolver) Resolve(ctx context.Context, tokenValue string) (oauth2.TokenInfo, error) {
	if v.TokenStore != nil {
		return v.TokenStore.GetByAccess(ctx, tokenValue)
	}
	return nil, nil
}
