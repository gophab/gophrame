package global

import (
	"log"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ERROR_BASE_PATH = 10011
	ErrorBasePath   = "无法获取根路径"
)

type CORS struct {
	Mode      string          `mapstructure:"mode" json:"mode" yaml:"mode"`
	Whitelist []CORSWhitelist `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`
}

type CORSWhitelist struct {
	AllowOrigin      string `mapstructure:"allow-origin" json:"allow-origin" yaml:"allow-origin"`
	AllowMethods     string `mapstructure:"allow-methods" json:"allow-methods" yaml:"allow-methods"`
	AllowHeaders     string `mapstructure:"allow-headers" json:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" json:"expose-headers" yaml:"expose-headers"`
	AllowCredentials bool   `mapstructure:"allow-credentials" json:"allow-credentials" yaml:"allow-credentials"`
}

var (
	BasePath string // 定义项目的根目录

	Debug bool = false

	EventDestroyPrefix = "Destroy_" //  程序退出时需要销毁的事件前缀
	ConfigKeyPrefix    = "Config_"  //  配置文件键值缓存时，键的前缀

	DateFormat = "2006-01-02 15:04:05" //  配置文件键值缓存时，键的前缀

	//gorm 数据库客户端，如果您操作数据库使用的是gorm，请取消以下注释，在 bootstrap>init 文件，进行初始化即可使用
	DB *gorm.DB // 全局gorm的客户端连接

	//gin engine
	Engine *gin.Engine

	//websocket
	WebsocketHub interface{}

	//casbin 全局操作指针
	Enforcer *casbin.SyncedEnforcer

	// CORS
	Cors *CORS = &CORS{
		Mode: "allow-all",
	}
)

func init() {
	// 1.初始化程序根目录
	if path, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(path, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = path
		}
		log.Println("Base application path: ", BasePath)
	} else {
		log.Fatal(ERROR_BASE_PATH, ErrorBasePath)
	}
}
