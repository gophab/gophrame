package server

import (
	"errors"

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
			c.Set("_CURRENT_USER_ID_", ti.GetUserID())
		}

		return ti, nil
	}
	return nil, errors.New(ErrorsNoToken)
}
