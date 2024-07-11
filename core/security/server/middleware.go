package server

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
)

const (
	ErrorsNoToken string = "no Token Verifier"
)

func ValidationBearerToken(c *gin.Context) (oauth2.TokenInfo, error) {
	// 1. 若EnableServer
	if theServer != nil {
		ti, err := theServer.ValidationBearerToken(c.Request)
		if err != nil {
			return nil, err
		}

		if ti != nil {
			var uid = strings.Split(ti.GetUserID(), "@")
			c.Set("_CURRENT_USER_ID_", uid[0])
			if len(uid) > 1 {
				c.Set("_CURRENT_TENANT_ID_", uid[1])
			}
		}

		return ti, nil
	}
	return nil, errors.New(ErrorsNoToken)
}
