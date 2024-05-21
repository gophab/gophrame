package config

type ConsulSetting struct {
	Enabled bool
}

var Setting *ConsulSetting = &ConsulSetting{
	Enabled: false,
}
