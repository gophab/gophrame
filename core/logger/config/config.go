package config

import "github.com/gophab/gophrame/core/config"

type LogSetting struct {
	LogName       string `json:"logName" yaml:"logName"`
	TextFormat    string `json:"textFormat" yaml:"textFormat"`
	TimePrecision string `json:"timePrecision" yaml:"timePrecision"`
	MaxSize       int    `json:"maxSize" yaml:"maxSize"`
	MaxBackups    int    `json:"maxBackups" yaml:"maxBackups"`
	MaxAge        int    `json:"maxAge" yaml:"maxAge"`
	Compress      bool   `json:"compress" yaml:"compress"`
}

var Setting *LogSetting = &LogSetting{}

func init() {
	config.RegisterConfig("log", Setting, "Log Settings")
}
