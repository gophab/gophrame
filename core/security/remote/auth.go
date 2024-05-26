package remote

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

type ClientAuthorizationSetting struct {
	AccessTokenURI string `json:"accessTokenUri" yaml:"accessTokenUri"`
	ClientId       string `json:"clientId" yaml:"clientId"`
	ClientSecret   string `json:"clientSecret" yaml:"clientSecret"`
}

var Setting *ClientAuthorizationSetting = &ClientAuthorizationSetting{}

type CheckInfo struct {
	Active    bool     `json:"active"`
	ExpiresIn int64    `json:"exp"`
	UserId    string   `json:"user_name"`
	ClientId  string   `json:"client_id"`
	Scope     []string `json:"scope"`
}

/**
 * POST /oauth/check_token?token={token}
 *
 *
 *
 */
func ValidationBearerToken(ctx *gin.Context) (oauth2.TokenInfo, error) {
	// 1. 获取当前Token
	tokenValue, err := SecurityUtil.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	if resp, err := http.Get(Setting.AccessTokenURI + "?token=" + tokenValue); err == nil && resp.StatusCode == 200 {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			checkInfo := &CheckInfo{}
			json.Json(string(body), &checkInfo)
			return &models.Token{
				Access:          tokenValue,
				ClientID:        checkInfo.ClientId,
				UserID:          checkInfo.UserId,
				Scope:           strings.Join(checkInfo.Scope, ","),
				AccessCreateAt:  time.Now(),
				AccessExpiresIn: time.Until(time.Unix(checkInfo.ExpiresIn, 0)),
			}, nil
		} else {
			logger.Error("Error reading token check info: ", err.Error())
			return nil, err
		}
	} else {
		logger.Error("Error remote calling: ", Setting.AccessTokenURI, err.Error())
		return nil, err
	}
}
