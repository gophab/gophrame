package jwt

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/golang-jwt/jwt"
)

type JwtSetting struct {
	Secret      string `json:"secret" yaml:"secret"`
	Method      string `json:"mehtod" yaml:"method"`
	OnlineUsers int    `json:"onlineUser" yaml:"onlineUsers"`
}

var (
	Setting *JwtSetting = &JwtSetting{
		Method: "HS256",
	}
)

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Scope       string `json:"scope"`
	RedirectURI string `json:"redirect_uri"`
	jwt.StandardClaims
}

// Valid claims verification
func (a *Claims) Valid() error {
	if time.Unix(a.ExpiresAt, 0).Before(time.Now()) {
		return errors.ErrInvalidAccessToken
	}
	return nil
}

func GenerateToken(claims *Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(signingMethod(), claims)
	if key, err := signingKey([]byte(Setting.Secret)); err == nil {
		return tokenClaims.SignedString(key)
	} else {
		return "", err
	}
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Setting.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func GetValueFromClaims(key string, claims jwt.Claims) string {
	v := reflect.ValueOf(claims)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)

			if fmt.Sprintf("%s", k.Interface()) == key {
				return fmt.Sprintf("%v", value.Interface())
			}
		}
	}
	return ""
}

func signingMethod() jwt.SigningMethod {
	var method jwt.SigningMethod = jwt.SigningMethodHS256

	if Setting.Method != "" {
		if v := jwt.GetSigningMethod(strings.ToUpper(Setting.Method)); v != nil {
			method = v
		}
	}

	return method
}

func signingKey(key []byte) (interface{}, error) {
	var result interface{}
	method := signingMethod()
	if isEs(method) {
		v, err := jwt.ParseECPrivateKeyFromPEM(key)
		if err != nil {
			return nil, err
		}
		result = v
	} else if isRsOrPS(method) {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(key)
		if err != nil {
			return nil, err
		}
		result = v
	} else if isHs(method) {
		result = key
	} else {
		return "", errors.New("unsupported sign method")
	}
	return result, nil
}

func isEs(method jwt.SigningMethod) bool {
	return strings.HasPrefix(method.Alg(), "ES")
}

func isRsOrPS(method jwt.SigningMethod) bool {
	isRs := strings.HasPrefix(method.Alg(), "RS")
	isPs := strings.HasPrefix(method.Alg(), "PS")
	return isRs || isPs
}

func isHs(method jwt.SigningMethod) bool {
	return strings.HasPrefix(method.Alg(), "HS")
}
