package config

import (
	CodeConfig "github.com/wjshen/gophrame/core/code/config"
	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"

	AliyunConfig "github.com/wjshen/gophrame/core/sms/aliyun/config"
	SmsCodeConfig "github.com/wjshen/gophrame/core/sms/code/config"
	QcloudConfig "github.com/wjshen/gophrame/core/sms/qcloud/config"
)

type SmsSetting struct {
	Enabled   bool `json:"enabled" yaml:"enabled"`
	Signature string
	Product   string
	Sender    struct {
		Aliyun *AliyunConfig.AliyunSetting `json:"aliyun" yaml:"aliyun"`
		Qcloud *QcloudConfig.QcloudSetting `json:"qcloud" yaml:"qcloud"`
	}
	Store *CodeConfig.CodeStoreSetting `json:"store" yaml:"store"`
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
	Store: SmsCodeConfig.Setting,
}

func init() {
	logger.Debug("Register SMS Config")
	config.RegisterConfig("sms", Setting, "SMS Settings")
}
