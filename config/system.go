package config

import (
	"time"

	c "github.com/wjshen/gophrame/core/config"
)

/**
 * 文件上传配置
 */
type FileUploadSetting struct {
	Size                 int      `json:"size"`
	UploadFileField      string   `json:"uploadFileField" yaml:"uploadFileField"`
	UploadFileSavePath   string   `json:"uploadFileSavePath" yaml:"uploadFileSavePath"`
	UploadFileReturnPath string   `json:"uploadFileReturnPath" yaml:"uploadFileReturnPath"`
	AllowMimeType        []string `json:"allowMimeType" yaml:"allowMimeType"`
}

var FileUpload *FileUploadSetting = &Config.FileUpload

/**
 * 服务配置：地址/端口
 */
type ServerSetting struct {
	BindAddr         string        `json:"bindAddr" yaml:"bindAddr"`
	Port             int           `json:"port"`
	ReadTimeout      time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout     time.Duration `json:"wirteTimeout" yaml:"writeTimeout"`
	AllowCrossDomain bool          `json:"allowCrossDomain" yaml:"allowCrossDomain"`
}

var Server = &Config.Server

/**
 * 全局配置
 */
type Configuration struct {
	Server     ServerSetting     `json:"server"`
	FileUpload FileUploadSetting `json:"fileUpload" yaml:"fileUpload"`
}

var Config *Configuration = &Configuration{}

func init() {
	c.RegisterConfig("ROOT", Config, "Default system configuration")
}
