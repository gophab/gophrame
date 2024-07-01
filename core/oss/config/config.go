package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"

	AliyunConfig "github.com/gophab/gophrame/core/oss/aliyun/config"
	QcloudConfig "github.com/gophab/gophrame/core/oss/qcloud/config"
)

type OssSetting struct {
	Enabled   bool `json:"enabled" yaml:"enabled"`
	Signature string
	Product   string
	Sender    struct {
		Aliyun *AliyunConfig.AliyunSetting `json:"aliyun" yaml:"aliyun"`
		Qcloud *QcloudConfig.QcloudSetting `json:"qcloud" yaml:"qcloud"`
	}
}

var Setting *OssSetting = &OssSetting{
	Enabled: false,
	Sender: struct {
		Aliyun *AliyunConfig.AliyunSetting `json:"aliyun" yaml:"aliyun"`
		Qcloud *QcloudConfig.QcloudSetting `json:"qcloud" yaml:"qcloud"`
	}{
		Aliyun: AliyunConfig.Setting,
		Qcloud: QcloudConfig.Setting,
	},
}

func init() {
	logger.Debug("Register SMS Config")
	config.RegisterConfig("sms", Setting, "SMS Settings")
}
