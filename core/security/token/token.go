package token

import (
	"github.com/gophab/gophrame/core/util"

	"github.com/go-oauth2/oauth2/v4/models"
)

type TokenRequest struct {
	ClientId          string            `json:"client_id"`
	Scope             string            `json:"scope"`
	GrantType         string            `json:"grant_type"`
	RequestParameters map[string]string `json:"request_parameters"`
}

type Authentication struct {
	UserId        string       `json:"user_id"`
	Request       TokenRequest `json:"request"`
	Authenticated bool         `json:"authenticated"`
}

func (a *Authentication) GetId() string {
	return util.MD5(a.UserId + a.Request.ClientId + a.Request.Scope)
}

type OAuth2Token struct {
	models.Token
}

// func (t *OAuth2Token) isExpired() bool {
// 	return t.AccessCreateAt.Add(t.AccessExpiresIn).Before(time.Now())
// }

// func (t *OAuth2Token) isRefreshExpired() bool {
// 	return t.RefreshCreateAt.Add(t.RefreshExpiresIn).Before(time.Now())
// }
