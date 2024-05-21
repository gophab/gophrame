package config

type DubboSetting struct {
	Enabled bool
}

var Setting *DubboSetting = &DubboSetting{
	Enabled: false,
}
