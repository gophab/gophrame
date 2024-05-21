package config

import (
	CodeConfig "github.com/wjshen/gophrame/core/code/config"

	AliyunConfig "github.com/wjshen/gophrame/core/sms/aliyun/config"
	"github.com/wjshen/gophrame/core/sms/code/config"
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
	Store: config.Setting,
}
