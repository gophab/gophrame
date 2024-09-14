package context

import (
	"sync"

	"github.com/gophab/gophrame/core/routine"
)

type GlobalContext struct {
	sync.RWMutex
	ContextVariables map[string]*routine.ThreadLocal[interface{}]
}

var globalContext = &GlobalContext{
	ContextVariables: make(map[string]*routine.ThreadLocal[interface{}]),
}

func (gc *GlobalContext) SetVariable(name string, v interface{}) {
	gc.Lock()
	defer gc.Unlock()

	variable, b := gc.ContextVariables[name]
	if !b {
		variable = routine.NewThreadLocal(v)

		gc.ContextVariables[name] = variable

	}
	variable.Set(v)
}

func (gc *GlobalContext) GetVariable(name string) interface{} {
	gc.RLock()
	defer gc.RUnlock()

	variable, b := gc.ContextVariables[name]
	if b {
		return variable.Get()
	}
	return nil
}

func (gc *GlobalContext) RemoveVariable(name string) {
	gc.Lock()
	defer gc.Unlock()

	variable, b := gc.ContextVariables[name]
	if b {
		variable.Remove()
		delete(gc.ContextVariables, name)
	}
}

func GetContextValue(name string) interface{} {
	return globalContext.GetVariable(name)
}

func SetContextValue(name string, v interface{}) {
	globalContext.SetVariable(name, v)
}

func RemoveContextValue(name string) {
	globalContext.RemoveVariable(name)
}
