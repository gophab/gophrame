package config

type ApplicationSetting struct {
	Debug bool `json:"debug"`

	PageSize  int
	PrefixUrl string

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

var Application *ApplicationSetting = &Config.Application

func init() {
	Application.ImageMaxSize = Application.ImageMaxSize * 1024 * 1024
}
