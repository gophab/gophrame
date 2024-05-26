package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type CaptchaSetting struct {
	Enabled      bool   `json:"enabled"`
	CaptchaId    string `json:"captchaId" yaml:"captchaId"`
	CaptchaValue string `json:"captchaValue" yaml:"captchaValue"`
	Length       int    `json:"length"`
}

var Setting *CaptchaSetting = &CaptchaSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Captcha Config")
	config.RegisterConfig("captcha", Setting, "Captcha Settings")
}
