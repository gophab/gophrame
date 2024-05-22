package config

import (
	"github.com/wjshen/gophrame/core/config"
)

type ApplicationSetting struct {
	Debug     bool   `json:"debug"`
	PageSize  int    `json:"pageSize" yaml:"pageSize"`
	PrefixUrl string `json:"prefixUrl" yaml:"prefixUrl"`

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var Application *ApplicationSetting = &ApplicationSetting{}

func init() {
	Application.ImageMaxSize = Application.ImageMaxSize * 1024 * 1024

	// ... 增加Application配置节点
	config.RegisterConfig("application", Application, "Application Settings")
}
