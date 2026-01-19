package config

import "github.com/gophab/gophrame/core/config"

type AlipaySetting struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	AppID     string `yaml:"appId" json:"appId"`
	MchID     string
	Key       string
	Secret    string
	NotifyURL string
	ReturnURL string
	IsProd    bool
	CertPath  string // 证书路径（微信支付需要）
}

var Setting = &AlipaySetting{
	Enabled: false,
}

func init() {
	config.RegisterConfig("payment.alipay", Setting, "Alipay Settings")
}
