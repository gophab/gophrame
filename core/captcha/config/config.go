package config

import (
	"time"

	CodeConfig "github.com/gophab/gophrame/core/code/config"
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type CaptchaSetting struct {
	Enabled      bool                         `json:"enabled"`
	CaptchaId    string                       `json:"captchaId" yaml:"captchaId"`
	CaptchaValue string                       `json:"captchaValue" yaml:"captchaValue"`
	Height       int                          `json:"height" yaml:"height"`
	Width        int                          `json:"width" yaml:"width"`
	Length       int                          `json:"length"`
	Store        *CodeConfig.CodeStoreSetting `json:"store"`
}

var Setting *CaptchaSetting = &CaptchaSetting{
	Enabled: false,
	Width:   240,
	Height:  80,
	Length:  6,
	Store: &CodeConfig.CodeStoreSetting{
		Enabled:         true,
		RequestInterval: time.Duration(10) * time.Second,
		ExpireIn:        time.Duration(5) * time.Minute,
	},
}

func init() {
	logger.Debug("Register Captcha Config")
	config.RegisterConfig("captcha", Setting, "Captcha Settings")
}
