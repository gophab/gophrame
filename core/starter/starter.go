package starter

import (
	"sort"

	"github.com/gophab/gophrame/core/logger"
)

type Func struct {
	f func()
	p int
}

var initializors = make([]*Func, 0)
var starters = make([]*Func, 0)
var terminaters = make([]*Func, 0)

func RegisterStarter(f func()) {
	RegisterStarterEx(f, int(0))
}

func RegisterStarterEx(f func(), p int) {
	starters = append(starters, &Func{f: f, p: p})
	sort.Slice(starters, func(i, j int) bool {
		return starters[i].p < starters[j].p
	})
}

func RegisterInitializor(f func()) {
	RegisterInitializorEx(f, int(0))
}

func RegisterInitializorEx(f func(), p int) {
	initializors = append(initializors, &Func{f: f, p: p})
	sort.Slice(initializors, func(i, j int) bool {
		return initializors[i].p < initializors[j].p
	})
}

func RegisterTerminater(f func()) {
	RegisterTerminaterEx(f, int(0))
}

func RegisterTerminaterEx(f func(), p int) {
	terminaters = append(terminaters, &Func{f: f, p: p})
	sort.Slice(terminaters, func(i, j int) bool {
		return terminaters[i].p < terminaters[j].p
	})
}

func Init() {
	logger.Info("Starting initializors...")
	for _, s := range initializors {
		s.f()
	}
}

func Start() {
	logger.Info("Starting starters...")
	for _, s := range starters {
		s.f()
	}
}

func Terminate() {
	logger.Info("Starting terminaters...")
	for _, s := range terminaters {
		s.f()
	}
}
