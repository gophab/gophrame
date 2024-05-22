package config

import (
	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"

	ServerConfig "github.com/wjshen/gophrame/core/security/server/config"
	TokenConfig "github.com/wjshen/gophrame/core/security/token/config"
)

type SecuritySetting struct {
	AuthMode     string `json:"authMode" yaml:"authMode"`
	AutoRegister bool   `json:"autoRegister" yaml:"autoRegister"`

	EmailAutoRegister  *bool `json:"emailAutoRegister" yaml:"emailAutoRegister"`
	MobileAutoRegister *bool `json:"mobileAutoRegister" yaml:"mobileAutoRegister"`
	SocialAutoRegister *bool `json:"socialAutoRegister" yaml:"socialAutoRegister"`

	// Server
	Server *ServerConfig.OAuth2ServerSetting `json:"server" yaml:"server"`

	// Token
	Token *TokenConfig.TokenSetting `json:"token" yaml:"token"`
}

var Setting *SecuritySetting = &SecuritySetting{
	AuthMode:     "local",
	AutoRegister: true,
	Server:       ServerConfig.Setting,
	Token:        TokenConfig.Setting,
}

func init() {
	logger.Debug("Register Security Config")
	config.RegisterConfig("security", Setting, "Security Settings")
}
