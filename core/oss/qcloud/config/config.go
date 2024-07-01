package config

type QcloudSetting struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	AppId     string `json:"appId" yaml:"appId"`
	AppKey    string `json:"appKey" yaml:"appKey"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Region    string `json:"region" yaml:"region"`
	BucketUrl string `json:"bucketUrl" yaml:"bucketUrl"`
	Path      string `json:"path" yaml:"path"`
}

var Setting *QcloudSetting = &QcloudSetting{
	Enabled: false,
}
