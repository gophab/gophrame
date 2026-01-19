package server

type ContextKey string

const AppIdContextKey = ContextKey("appId")
const AuthorizationCodeKey = ContextKey("authorizationCode")
