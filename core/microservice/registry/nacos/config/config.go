package config

type NacosSetting struct {
	Enabled bool
}

var Setting *NacosSetting = &NacosSetting{
	Enabled: false,
}
