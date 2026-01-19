package service

import (
	"strings"
	"sync"
)

type EntityGetter func(id any) any

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

func (h *EntityHelper) GetEntity(entity string, id any) any {
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

func GetEntity(entity string, id any) any {
	return entityHelper.GetEntity(entity, id)
}

func GetEntityAs[T any](entity string, id any) *T {
	return entityHelper.GetEntity(entity, id).(*T)
}
