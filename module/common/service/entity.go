package service

import (
	"strings"
	"sync"
)

type EntityGetter func(id interface{}) interface{}

type EntityHelper struct {
	mutex   sync.RWMutex
	Getters map[string]EntityGetter
}

var entityHelper = &EntityHelper{
	mutex:   sync.RWMutex{},
	Getters: make(map[string]EntityGetter),
}

func (h *EntityHelper) RegisterEntity(entity string, getter EntityGetter) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Getters[strings.ToLower(entity)] = getter
}

func (h *EntityHelper) GetEntity(entity string, id interface{}) interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if getter, b := h.Getters[strings.ToLower(entity)]; b {
		return getter(id)
	}
	return nil
}

func RegisterEntity(entity string, getter EntityGetter) {
	entityHelper.RegisterEntity(entity, getter)
}

func GetEntity(entity string, id interface{}) interface{} {
	return entityHelper.GetEntity(entity, id)
}

func GetEntityAs[T any](entity string, id interface{}) *T {
	return entityHelper.GetEntity(entity, id).(*T)
}
