package util

func IntAddr(i int) *int {
	result := new(int)
	*result = i
	return result
}

func Int32Addr(i int32) *int32 {
	result := new(int32)
	*result = i
	return result
}

func Int64Addr(i int64) *int64 {
	result := new(int64)
	*result = i
	return result
}

func IntValue(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func DefaultIntValue(i *int, defaultValue int) int {
	if i != nil {
		return *i
	}
	return defaultValue
}

func Int64Value(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func DefaultInt64Value(i *int64, defaultValue int64) int64 {
	if i != nil {
		return *i
	}
	return defaultValue
}

func Int32Value(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

func DefaultInt32Value(i *int32, defaultValue int32) int32 {
	if i != nil {
		return *i
	}
	return defaultValue
}
