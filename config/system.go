package config

import (
	"time"

	"github.com/gophab/gophrame/core/config"

	_ "github.com/gophab/gophrame/core/captcha/config"
	_ "github.com/gophab/gophrame/core/casbin/config"
	_ "github.com/gophab/gophrame/core/database/config"
	_ "github.com/gophab/gophrame/core/email/config"
	_ "github.com/gophab/gophrame/core/logger/config"
	_ "github.com/gophab/gophrame/core/microservice/config"
	_ "github.com/gophab/gophrame/core/rabbitmq/config"
	_ "github.com/gophab/gophrame/core/redis/config"
	_ "github.com/gophab/gophrame/core/security/config"
	_ "github.com/gophab/gophrame/core/sms/config"
	_ "github.com/gophab/gophrame/core/snowflake/config"
	_ "github.com/gophab/gophrame/core/social/config"
	_ "github.com/gophab/gophrame/core/websocket/config"
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
	config.RegisterConfig("ROOT", Config, "Default system configuration")
}
