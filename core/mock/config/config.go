package config

import (
	"github.com/gophab/gophrame/core/config"
)

type MockSetting struct {
	Enabled bool
	Apis    map[string]string
}

var Setting *MockSetting = &MockSetting{
	Enabled: false,
}

func init() {
	config.RegisterConfig("mock", Setting, "Mock Settings")
}
