package config

import "github.com/gophab/gophrame/core/config"

type TCloudSetting struct {
	Enabled bool
}

var Setting = &TCloudSetting{
	Enabled: false,
}

func init() {
	config.RegisterConfig("identity.tcloud", Setting, "Tencent Cloud Identify Settings")
}
