package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"

	AliyunConfig "github.com/gophab/gophrame/core/sms/aliyun/config"
	QcloudConfig "github.com/gophab/gophrame/core/sms/qcloud/config"
)

type SmsSetting struct {
	Enabled   bool `json:"enabled" yaml:"enabled"`
	Signature string
	Product   string
	Sender    struct {
		Aliyun *AliyunConfig.AliyunSetting `json:"aliyun" yaml:"aliyun"`
		Qcloud *QcloudConfig.QcloudSetting `json:"qcloud" yaml:"qcloud"`
	}
}

var Setting *SmsSetting = &SmsSetting{
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
