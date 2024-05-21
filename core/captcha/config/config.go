package config

type CaptchaSetting struct {
	Enabled      bool   `json:"enabled"`
	CaptchaId    string `json:"captchaId" yaml:"captchaId"`
	CaptchaValue string `json:"captchaValue" yaml:"captchaValue"`
	Length       int    `json:"length"`
}

var Setting *CaptchaSetting = &CaptchaSetting{
	Enabled: false,
}
