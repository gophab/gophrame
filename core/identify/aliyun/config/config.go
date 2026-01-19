package config

import "github.com/gophab/gophrame/core/config"

type AliyunSetting struct {
	Enabled        bool   `yam:"enabled" json:"enabled"`
	Base           string `yaml:"base" json:"base"`
	Url            string `yaml:"url" json:"url"`
	Proxy          string `yaml:"proxy" json:"proxy"`
	TwoFactorUrl   string `yaml:"twoFactorUrl" json:"twoFactorUrl"`
	ThreeFactorUrl string `yaml:"threeFactorUrl" json:"threeFactorUrl"`
	AppCode        string `yaml:"appCode" json:"appCode"`
}

var Setting = &AliyunSetting{
	Enabled: false,
}

func init() {
	config.RegisterConfig("identify.aliyun", Setting, "Aliyun Identify Settings")
}
