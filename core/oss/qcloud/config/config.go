package config

type QcloudSetting struct {
	Enabled   bool              `json:"enabled" yaml:"enabled"`
	AppId     string            `json:"appId" yaml:"appId"`
	AppKey    string            `json:"appKey" yaml:"appKey"`
	Templates map[string]string `json:"templates" yaml:"templates"`
}

var Setting *QcloudSetting = &QcloudSetting{
	Enabled: false,
}
