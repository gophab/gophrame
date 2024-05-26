package module

import (
	"sort"
	"sync"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

type IModule interface {
	// Call when register
	Register()

	// Call when initializing, after init config
	Init()

	// Call when start, after root router initialized
	Start()

	// Call before exit
	Terminate()
}

const (
	STATUS = iota
	STATUS_REGISTERED
	STATUS_INITIALIZED
	STATUS_STARTED
	STATUS_TERMINATED
)

type Module struct {
	IModule
	Name        string
	Description string
	Initializor func(m *Module)
	Starter     func(m *Module)
	Terminater  func(m *Module)
	Priority    int
	Status      int
}

func (m *Module) Terminate() {
	if m.Terminater != nil {
		logger.Debug("[MODULE] Terminating module: ", m.Name)
		m.Terminater(m)
	}
}

func (m *Module) Start() {
	if m.Starter != nil {
		logger.Debug("[MODULE] Starting module: ", m.Name)
		m.Starter(m)
	}
}

func (m *Module) Init() {
	if m.Initializor != nil {
		logger.Debug("[MODULE] Initializing module: ", m.Name)
		m.Initializor(m)
	}
}

func (*Module) Register() {}

var modules = make([]*Module, 0)

func RegisterModule(mod *Module) {
	modules = append(modules, mod)
	mod.Register()
	mod.Status = STATUS_REGISTERED
}

var (
	onceInit, onceStart sync.Once
)

func Init() {
	onceInit.Do(func() {
		sort.Slice(modules, func(i, j int) bool {
			return modules[i].Priority < modules[j].Priority
		})
		for _, mod := range modules {
			mod.Init()
			mod.Status = STATUS_INITIALIZED
		}
	})
}

func Start() {
	onceStart.Do(func() {
		for _, mod := range modules {
			mod.Start()
			mod.Status = STATUS_STARTED
		}
	})
}

func Terminate() {
	onceStart.Do(func() {
		for _, mod := range modules {
			mod.Terminate()
			mod.Status = STATUS_TERMINATED
		}
	})
}

func init() {
	starter.RegisterInitializorEx(Init, 0x0FFFFFFF)
	starter.RegisterStarterEx(Start, 0x0FFFFFFF)
	starter.RegisterTerminaterEx(Terminate, -0x0FFFFFFF)
}
