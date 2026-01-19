package context

import (
	"sync"

	"github.com/gophab/gophrame/core/routine"
)

type GlobalContext struct {
	sync.RWMutex
	ContextVariables map[string]*routine.ThreadLocal[any]
}

var globalContext = &GlobalContext{
	ContextVariables: make(map[string]*routine.ThreadLocal[any]),
}

func (gc *GlobalContext) SetVariable(name string, v any) {
	gc.Lock()
	defer gc.Unlock()

	variable, b := gc.ContextVariables[name]
	if !b {
		variable = routine.NewThreadLocal(v)

		gc.ContextVariables[name] = variable

	}
	variable.Set(v)
}

func (gc *GlobalContext) GetVariable(name string) any {
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
	}
}

func GetContextValue(name string) any {
	return globalContext.GetVariable(name)
}

func SetContextValue(name string, v any) {
	globalContext.SetVariable(name, v)
}

func RemoveContextValue(name string) {
	globalContext.RemoveVariable(name)
}
