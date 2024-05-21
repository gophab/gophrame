package eventbus

import (
	"strings"
	"sync"

	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/errors"
)

// 公共消息总线
var theEventbus *eventbus = CreateEventbus()

func init() {
	inject.InjectValue("eventbus", theEventbus)
}

func RegisterEventListener(key string, keyFunc func(args ...interface{})) bool {
	return theEventbus.RegisterEventListener(key, keyFunc)
}

func RemoveEventListener(key string, keyFunc func(args ...interface{})) {
	theEventbus.RegisterEventListener(key, keyFunc)
}

/**
 * 同步分发消息
 */
func PublishEvent(key string, args ...interface{}) {
	theEventbus.PublishEvent(key, args...)
}

/**
 *	异步分发消息
 */
func DispatchEvent(key string, args ...interface{}) {
	theEventbus.DispatchEvent(key, args...)
}

/**
 * 模糊匹配分发消息
 */
func FuzzyPublishEvent(keyPre string, args ...interface{}) {
	theEventbus.FuzzyPublishEvent(keyPre, args...)
}

/**
 * 模糊匹配分发消息
 */
func FuzzyDispatchEvent(keyPre string, args ...interface{}) {
	theEventbus.FuzzyDispatchEvent(keyPre, args...)
}

// 创建一个事件管理工厂
func CreateEventbus() *eventbus {
	return new(eventbus)
}

// 定义一个事件管理结构体
type eventbus struct {
	// 定义一个全局事件存储变量，本模块只负责存储 键 => 函数 ， 相对容器来说功能稍弱，但是调用更加简单、方便、快捷
	eventListeners sync.Map
}

// 1.注册事件
func (e *eventbus) RegisterEventListener(key string, keyFunc func(args ...interface{})) bool {
	//判断key下是否已有事件
	if queue, exists := e.GetEventListeners(key); !exists {
		e.eventListeners.Store(key, []func(args ...interface{}){keyFunc})
		return true
	} else {
		e.eventListeners.Store(key, append(queue, keyFunc))
	}
	return false
}

// 2.获取事件
func (e *eventbus) GetEventListeners(key string) ([]func(args ...interface{}), bool) {
	if queue, exists := e.eventListeners.Load(key); exists {
		return queue.([]func(args ...interface{})), exists
	}
	return nil, false
}

// 3.执行事件
func (e *eventbus) PublishEvent(key string, args ...interface{}) {
	if queue, exists := e.GetEventListeners(key); exists {
		for _, fn := range queue {
			fn(args...)
		}
	} else {
		logger.Error(errors.ERROR_FUNC_EVENT_NOT_REGISTER, ", 键名：", key)
	}
}

// 3.执行事件
func (e *eventbus) DispatchEvent(key string, args ...interface{}) {
	if queue, exists := e.GetEventListeners(key); exists {
		for _, fn := range queue {
			go func(cb func(args ...interface{})) {
				cb(args...)
			}(fn)
		}
	} else {
		logger.Error(errors.ERROR_FUNC_EVENT_NOT_REGISTER, ", 键名：", key)
	}
}

// 4.删除事件
func (e *eventbus) RemoveEventListeners(key string) {
	e.eventListeners.Delete(key)
}

// 4.删除事件
func (e *eventbus) RemoveEventListener(key string, keyFunc func(args ...interface{})) {
	if queue, exists := e.GetEventListeners(key); exists {
		var j = 0
		for _, v := range queue {
			if &keyFunc != &v {
				queue[j] = v
				j++
			}
		}
		e.eventListeners.Store(key, queue[:j])
	}
}

// 5.根据键的前缀，模糊调用. 使用请谨慎.
func (e *eventbus) FuzzyPublishEvent(keyPre string, args ...interface{}) {
	e.eventListeners.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.PublishEvent(keyName, args...)
			}
		}
		return true
	})
}

// 6.根据键的前缀，模糊调用. 使用请谨慎.
func (e *eventbus) FuzzyDispatchEvent(keyPre string, args ...interface{}) {
	e.eventListeners.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.DispatchEvent(keyName, args...)
			}
		}
		return true
	})
}
