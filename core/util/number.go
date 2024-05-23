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
