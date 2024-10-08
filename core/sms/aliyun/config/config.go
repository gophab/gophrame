package config

type AliyunSetting struct {
	Enabled         bool              `json:"enabled" yaml:"enabled"`
	AccessKeyId     string            `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string            `json:"accessKeySecret" yaml:"accessKeySecret"`
	Signature       string            `json:"signature" yaml:"signature"`
	SignatureTC     string            `json:"signatureTC" yaml:"signatureTC"`
	SignatureEN     string            `json:"signatureEN" yaml:"signatureEN"`
	Product         string            `json:"product" yaml:"product"`
	Templates       map[string]string `json:"templates" yaml:"templates"`
}

var Setting *AliyunSetting = &AliyunSetting{
	Enabled: false,
}
