package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type AliyunSetting struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	UseCname        bool   `json:"useCname" yaml:"useCname"`
	Bucket          string `json:"bucket" yaml:"bucket"`
	Region          string `json:"region" yaml:"region"`
	BucketUrl       string `json:"bucketUrl" yaml:"bucketUrl"`
	Path            string `json:"path" yaml:"path"`
}

var Setting *AliyunSetting = &AliyunSetting{
	Enabled:  false,
	UseCname: false,
	Path:     "",
}

func init() {
	logger.Debug("Register Oss Config - Aliyun")
	config.RegisterConfig("oss.aliyun", Setting, "Aliyun Settings")
}
