package config

import "github.com/gophab/gophrame/core/config"

type WxpaySetting struct {
	Enabled             bool   `yaml:"enabled" json:"enabled"`
	AppID               string `yaml:"appId" json:"appId"`
	MchID               string `yaml:"mchId" json:"mchId"`
	APIv3Key            string `yaml:"apiV3Key" json:"apiV3Key"`
	CertificateSerialNo string `yaml:"certificateSerialNo" json:"certificateSerialNo"`
	PrivateKeyFilePath  string `yaml:"privateKeyFilePath" json:"privateKeyFilePath"` // 证书路径（微信支付需要）
	NotifyURL           string `yaml:"notifyUrl" json:"notifyUrl"`
	ReturnURL           string `yaml:"returnUrl" json:"returnUrl"`
	IsProd              bool   `yaml:"isProd" json:"isProd"`
}

var Setting = &WxpaySetting{
	Enabled: false,
}

func init() {
	config.RegisterConfig("payment.wxpay", Setting, "Wxpay Settings")
}
