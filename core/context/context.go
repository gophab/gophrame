package context

import "github.com/gophab/gophrame/core/routine"

type GlobalContext struct {
	ContextVariables map[string]*routine.ThreadLocal[interface{}]
}

var globalContext = &GlobalContext{
	ContextVariables: make(map[string]*routine.ThreadLocal[interface{}]),
}

func (gc *GlobalContext) SetVariable(name string, v interface{}) {
	variable, b := gc.ContextVariables[name]
	if !b {
		variable = routine.NewThreadLocal(v)
		gc.ContextVariables[name] = variable
	}
	variable.Set(v)
}

func (gc *GlobalContext) GetVariable(name string) interface{} {
	variable, b := gc.ContextVariables[name]
	if b {
		return variable.Get()
	}
	return nil
}

func (gc *GlobalContext) RemoveVariable(name string) {
	variable, b := gc.ContextVariables[name]
	if b {
		variable.Remove()
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
