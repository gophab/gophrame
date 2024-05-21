package config

type EurekaSetting struct {
	Enabled     bool
	ServiceUrls []string
	Username    string
	Password    string
}

var Setting *EurekaSetting = &EurekaSetting{
	Enabled: false,
}
