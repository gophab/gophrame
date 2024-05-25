package global

import "time"

var ()

var global_values = make(map[string]interface{})

func Set(name string, value interface{}) {
	global_values[name] = value
}

func Sets(vars map[string]interface{}) {
	for key, value := range vars {
		global_values[key] = value
	}
}

func Var(name string) interface{} {
	if v, b := global_values[name]; b {
		return v
	} else {
		return nil
	}
}

func IntVar(name string) int {
	return VarValue(name, int(0))
}

func IntVarDefault(name string, v int) int {
	return VarValue(name, int(v))
}

func Int64Var(name string) int64 {
	return VarValue(name, int64(0))
}

func Int64VarDefault(name string, v int64) int64 {
	return VarValue(name, int64(v))
}

func StringVar(name string) string {
	return VarValue(name, "")
}

func StringVarDefault(name string, v string) string {
	return VarValue(name, v)
}

func DurationVar(name string) time.Duration {
	if value := Var(name); value != nil {
		return parserDuration(value)
	}
	return time.Duration(0)
}

func DurationVarDefault(name string, v interface{}) time.Duration {
	if value := Var(name); value != nil {
		return parserDuration(value)
	}
	return parserDuration(v)
}

func parserDuration(v interface{}) time.Duration {
	switch v := v.(type) {
	case int:
		return time.Second * time.Duration(v)
	case int32:
		return time.Second * time.Duration(v)
	case int64:
		return time.Second * time.Duration(v)
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return time.Duration(0)
}

func VarValue[T any](name string, t T) T {
	if v := Var(name); v != nil {
		return v.(T)
	} else {
		return t
	}
}
