package config

import (
	"time"

	c "github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/json"
	"github.com/wjshen/gophrame/core/logger"

	CaptchaConfig "github.com/wjshen/gophrame/core/captcha/config"
	CasbinConfig "github.com/wjshen/gophrame/core/casbin/config"
	DatabaseConfig "github.com/wjshen/gophrame/core/database/config"
	EmailConfig "github.com/wjshen/gophrame/core/email/config"
	LoggerConfig "github.com/wjshen/gophrame/core/logger/config"
	RabbitMQConfig "github.com/wjshen/gophrame/core/rabbitmq/config"
	RedisConfig "github.com/wjshen/gophrame/core/redis/config"
	SecurityConfig "github.com/wjshen/gophrame/core/security/config"
	SmsConfig "github.com/wjshen/gophrame/core/sms/config"
	SnowflakeConfig "github.com/wjshen/gophrame/core/snowflake/config"
	SocialConfig "github.com/wjshen/gophrame/core/social/config"
	WebsocketConfig "github.com/wjshen/gophrame/core/websocket/config"
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
	// ... 增加Application配置节点
	Application ApplicationSetting `json:"application"`
	Server      ServerSetting      `json:"server"`
	FileUpload  FileUploadSetting  `json:"fileUpload" yaml:"fileUpload"`

	Security  *SecurityConfig.SecuritySetting   `json:"security"`
	Database  *DatabaseConfig.DatabaseSetting   `json:"database"`
	Redis     *RedisConfig.RedisSetting         `json:"redis"`
	Log       *LoggerConfig.LogSetting          `json:"log"`
	SnowFlake *SnowflakeConfig.SnowFlakeSetting `json:"snowflake"`
	Captcha   *CaptchaConfig.CaptchaSetting     `json:"captcha"`
	Sms       *SmsConfig.SmsSetting             `json:"sms" yaml:"sms"`
	Email     *EmailConfig.EmailSetting         `json:"email" yaml:"email"`
	Casbin    *CasbinConfig.CasbinSetting       `json:"casbin"`
	RabbitMQ  *RabbitMQConfig.RabbitMQSetting   `json:"rabbitmq" yaml:"rabbitmq"`
	Websocket *WebsocketConfig.WebsocketSetting `json:"websocket"`
	Social    *SocialConfig.SocialSetting       `json:"social"`
}

var Config *Configuration = &Configuration{
	Database:  DatabaseConfig.Setting,
	SnowFlake: SnowflakeConfig.Setting,
	Casbin:    CasbinConfig.Setting,
	Redis:     RedisConfig.Setting,
	Websocket: WebsocketConfig.Setting,
	Log:       LoggerConfig.Setting,
	Captcha:   CaptchaConfig.Setting,
	Sms:       SmsConfig.Setting,
	Email:     EmailConfig.Setting,
	RabbitMQ:  RabbitMQConfig.Setting,
	Security:  SecurityConfig.Setting,
	Social:    SocialConfig.Setting,
}

var ConfigYml c.IYmlConfig = c.ConfigYml

func init() {
	c.InitConfig(&Config)
	logger.Debug("Load application configuration: ", json.String(Config))
}
