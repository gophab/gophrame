package config

type AliyunSetting struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
	Signature       string `json:"signature" yaml:"signature"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	UseCname        bool   `json:"useCname" yaml:"useCname"`
	Bucket          string `json:"bucket" yaml:"bucket"`
	BucketUrl       string `json:"bucketUrl" yaml:"bucketUrl"`
	Path            string `json:"path" yaml:"path"`
}

var Setting *AliyunSetting = &AliyunSetting{
	Enabled:  false,
	UseCname: false,
	Path:     "",
}
