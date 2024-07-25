package eventbus

import (
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/errors"
)

type EventListener func(event string, args ...interface{})

// 公共消息总线
var theEventbus *eventbus = CreateEventbus()

func init() {
	inject.InjectValue("eventbus", theEventbus)
}

func RegisterEventListener(event string, listener EventListener) bool {
	return theEventbus.RegisterEventListener(event, listener)
}

func RemoveEventListener(event string, listener EventListener) {
	theEventbus.RegisterEventListener(event, listener)
}

/**
 * 同步分发消息
 */
func PublishEvent(event string, args ...interface{}) {
	theEventbus.PublishEvent(event, args...)
}

/**
 *	异步分发消息
 */
func DispatchEvent(event string, args ...interface{}) {
	theEventbus.DispatchEvent(event, args...)
}

/**
 * 模糊匹配分发消息
 */
func FuzzyPublishEvent(eventPre string, args ...interface{}) {
	theEventbus.FuzzyPublishEvent(eventPre, args...)
}

/**
 * 模糊匹配分发消息
 */
func FuzzyDispatchEvent(eventPre string, args ...interface{}) {
	theEventbus.FuzzyDispatchEvent(eventPre, args...)
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
func (e *eventbus) RegisterEventListener(event string, listener EventListener) bool {
	//判断key下是否已有事件
	if queue, exists := e.GetEventListeners(event); !exists {
		e.eventListeners.Store(event, []EventListener{listener})
		return true
	} else {
		e.eventListeners.Store(event, append(queue, listener))
	}
	return false
}

// 2.获取事件
func (e *eventbus) GetEventListeners(event string) ([]EventListener, bool) {
	if queue, exists := e.eventListeners.Load(event); exists {
		return queue.([]EventListener), exists
	}
	return nil, false
}

// 3.执行事件
func (e *eventbus) PublishEvent(event string, args ...interface{}) {
	if queue, exists := e.GetEventListeners(event); exists {
		for _, fn := range queue {
			fn(event, args...)
		}
	} else {
		logger.Warn(errors.ERROR_FUNC_EVENT_NOT_REGISTER, ", 无效键名：", event)
	}
}

// 3.执行事件
func (e *eventbus) DispatchEvent(event string, args ...interface{}) {
	if queue, exists := e.GetEventListeners(event); exists {
		for _, fn := range queue {
			go func(cb EventListener) {
				cb(event, args...)
			}(fn)
		}
	} else {
		logger.Error(errors.ERROR_FUNC_EVENT_NOT_REGISTER, ", 无效键名：", event)
	}
}

// 4.删除事件
func (e *eventbus) RemoveEventListeners(event string) {
	e.eventListeners.Delete(event)
}

// 4.删除事件
func (e *eventbus) RemoveEventListener(event string, listener EventListener) {
	if queue, exists := e.GetEventListeners(event); exists {
		var j = 0
		for _, v := range queue {
			if &listener != &v {
				queue[j] = v
				j++
			}
		}
		e.eventListeners.Store(event, queue[:j])
	}
}

// 5.根据键的前缀，模糊调用. 使用请谨慎.
func (e *eventbus) FuzzyPublishEvent(eventPre string, args ...interface{}) {
	e.eventListeners.Range(func(eventKey, value interface{}) bool {
		if event, ok := eventKey.(string); ok {
			if strings.HasPrefix(event, eventPre) {
				e.PublishEvent(event, args...)
			}
		}
		return true
	})
}

// 6.根据键的前缀，模糊调用. 使用请谨慎.
func (e *eventbus) FuzzyDispatchEvent(eventPre string, args ...interface{}) {
	e.eventListeners.Range(func(eventKey, value interface{}) bool {
		if event, ok := eventKey.(string); ok {
			if strings.HasPrefix(event, eventPre) {
				e.DispatchEvent(event, args...)
			}
		}
		return true
	})
}
