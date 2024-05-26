package plugin

import (
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

const (
	STATUS = iota
	STATUS_REGISTERED
	STATUS_INITIALIZED
	STATUS_STARTED
	STATUS_TERMINATED
)

// Plugin 插件接口
type IPlugin interface {
	// Register 注册路由
	Register(group *gin.RouterGroup)

	// RouterPath 用户返回注册路由
	RouterPath() string
}

type Plugin struct {
	IPlugin
	Name        string
	Description string
	UUID        string // universal identify, used by Plugin Market
	Priority    int
	Status      int
	Enabled     bool
}

func (*Plugin) Terminate() {}

func (*Plugin) Start() {}

func (*Plugin) Init() {}

func (*Plugin) Register() {}

var plugins = make([]*Plugin, 0)

func RegisterPlugin(plugin *Plugin) {
	plugins = append(plugins, plugin)
	plugin.Register()
	plugin.Status = STATUS_REGISTERED
}

var (
	onceInit, onceStart sync.Once
)

func Init() {
	onceInit.Do(func() {
		sort.Slice(plugins, func(i, j int) bool {
			return plugins[i].Priority < plugins[j].Priority
		})
		for _, plugin := range plugins {
			plugin.Init()
			plugin.Status = STATUS_INITIALIZED

			// Plugin Engine
			if plugin.Enabled {
				if err := engine.InitPlugin(plugin); err != nil {
					// 引擎初始化
					logger.Error("[PLUGIN] Initializing plugin [", plugin.Name, "] error: ", err.Error())
				}
			}
		}
	})
}

func Start() {
	onceStart.Do(func() {
		for _, plugin := range plugins {
			// Plugin Engine
			if plugin.Enabled && plugin.Status == STATUS_INITIALIZED {
				plugin.Start()
				plugin.Status = STATUS_STARTED

				if err := engine.StartPlugin(plugin); err != nil {
					// 引擎初始化
					logger.Error("[PLUGIN] Starting plugin [", plugin.Name, "] error: ", err.Error())
				}
			}
		}
	})
}

func Terminate() {
	onceStart.Do(func() {
		for _, plugin := range plugins {
			// Plugin Engine
			if plugin.Enabled && plugin.Status == STATUS_STARTED {
				plugin.Terminate()
				plugin.Status = STATUS_TERMINATED
				if err := engine.StartPlugin(plugin); err != nil {
					// 引擎初始化
					logger.Error("[PLUGIN] Starting plugin [", plugin.Name, "] error: ", err.Error())
				}
			}
		}
	})
}

func init() {
	starter.RegisterInitializorEx(Init, 0x1FFFFFFF)
	starter.RegisterStarterEx(Start, 0x1FFFFFFF)
	starter.RegisterTerminaterEx(Terminate, -0x1FFFFFFF)
}
