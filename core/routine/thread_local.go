package routine

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/timandy/routine"
)

type LocalStorage[T any] map[string]T

func (ls *LocalStorage[T]) Get() (r T) {
	// 1. get current routing id
	gid := routine.Goid()
	var vmap map[string]T = *ls
	if value, b := vmap[strconv.FormatInt(gid, 10)]; b {
		return value
	}
	return r
}

func (ls *LocalStorage[T]) Set(value T) {
	gid := routine.Goid()
	var vmap map[string]T = *ls
	vmap[strconv.FormatInt(gid, 10)] = value
}

func (ls *LocalStorage[T]) Remove() {
	gid := routine.Goid()
	var vmap map[string]T = *ls
	delete(vmap, strconv.FormatInt(gid, 10))
}

func (ls *LocalStorage[T]) Clear() {
	var vmap map[string]T = *ls
	for k := range vmap {
		delete(vmap, k)
	}
}

type IThreadLocal interface {
	Clear()
	Index() int
}

type ThreadLocal[T any] struct {
	LocalStorage[T]
	index int
}

func (tl *ThreadLocal[T]) Index() int {
	return tl.index
}

var threadLocalIndex int32 = -1

func nextThreadLocalIndex() int {
	index := atomic.AddInt32(&threadLocalIndex, 1)
	if index < 0 {
		atomic.AddInt32(&threadLocalIndex, -1)
		panic("too many thread-local indexed variables")
	}
	return int(index)
}

func NewThreadLocal[T any](value T) *ThreadLocal[T] {
	var result = &ThreadLocal[T]{LocalStorage: make(map[string]T), index: nextThreadLocalIndex()}
	threadLocalManager.Add(result)
	return result
}

func Clear() {
	threadLocalManager.Clear()
}

type ThreadLocalManager struct {
	threadLocals map[int]IThreadLocal
	mutex        sync.Mutex
}

var threadLocalManager = &ThreadLocalManager{threadLocals: make(map[int]IThreadLocal), mutex: sync.Mutex{}}

func (m *ThreadLocalManager) Add(threadLocal IThreadLocal) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.threadLocals[threadLocal.Index()] = threadLocal
}

func (m *ThreadLocalManager) Remove(index int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.threadLocals, index)
}

func (m *ThreadLocalManager) RemoveThreadLocal(threadLocal IThreadLocal) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.threadLocals, threadLocal.Index())
}

func (m *ThreadLocalManager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, threadLocal := range m.threadLocals {
		threadLocal.Clear()
	}
	for k := range m.threadLocals {
		delete(m.threadLocals, k)
	}
}
