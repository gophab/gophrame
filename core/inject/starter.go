package inject

import (
	"github.com/gophab/gophrame/core/logger"
)

type AfterInitialize interface {
	AfterInitialize()
}

var values = make(map[string]any)
var graph Graph

// 初始化依赖注入
func Init() {
	for _, value := range values {
		_ = graph.Provide(&Object{Value: value})
	}

	if err := graph.Populate(); err != nil {
		logger.Error("初始化依赖注入发生错误：", err.Error())
	}
}

func InjectValue(key string, value any) {
	_ = graph.Provide(&Object{Name: key, Value: value})
	_ = graph.Populate()
	values[key] = value

	if iface, ok := value.(AfterInitialize); ok {
		iface.AfterInitialize()
	}
}

func InjectValue_(key string, value any) {
	values[key] = value
}

func GetValue(key string) any {
	return values[key]
}
