package plugin

type EntryPoint struct {
	NameSpace string
	Code      string
	Entry     func(args ...any)
}

type CallbackPoint struct {
	NameSpace string
	Code      string
	Callbacks []func(args ...any)
}

var EntryPoints = make(map[string]EntryPoint)
var CallbackPoints = make(map[string]CallbackPoint)

func RegisterEntryPoint(entryPoint *EntryPoint) {
	EntryPoints[entryPoint.NameSpace+":"+entryPoint.Code] = *entryPoint
}

func RegisterCallbackPoint(callback *CallbackPoint) {
	CallbackPoints[callback.NameSpace+":"+callback.Code] = *callback
}

// call from internal
func Callback(namespace, code string, args ...any) {
	if cp, b := CallbackPoints[namespace+":"+code]; b {
		for _, callback := range cp.Callbacks {
			callback(args)
		}
	}
}

type Engine struct {
}

var engine = &Engine{}

func (e *Engine) InitPlugin(plugin *Plugin) error {
	return nil
}

func (*Engine) StartPlugin(plugin *Plugin) error {
	return nil
}

func (*Engine) RegisterCallback(namespace, code string, f func(args ...any)) {
	if cp, b := CallbackPoints[namespace+":"+code]; b {
		cp.Callbacks = append(cp.Callbacks, f)
	}
}

// Call from plugin
func (*Engine) CallEntryPoint(namespace, code string, args any) {
	if entry, b := EntryPoints[namespace+":"+code]; b {
		entry.Entry(args)
	}
}
