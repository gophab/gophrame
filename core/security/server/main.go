package server

import "github.com/go-oauth2/oauth2/v4/server"

func RegisterClientScopeHandler(handler server.ClientScopeHandler) {
	if theServer != nil {
		theServer.RegisterClientScopeHandler(handler)
	}
}

func RegisterClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	if theServer != nil {
		theServer.RegisterClientAuthorizedHandler(handler)
	}
}

func RegisterPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler) {
	if theServer != nil {
		theServer.RegisterPasswordAuthorizationHandler(handler)
	}
}

func RegisterUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	if theServer != nil {
		theServer.RegisterUserAuthorizationHandler(handler)
	}
}
