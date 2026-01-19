package util

func If[T any](b bool, values ...T) T {
	if b {
		if len(values) > 0 {
			return values[0]
		}
	} else {
		if len(values) > 1 {
			return values[1]
		}
	}

	var empty T
	return []T{empty}[0]
}

func IfNot[T any](b bool, values ...T) T {
	if !b {
		if len(values) > 0 {
			return values[0]
		}
	} else {
		if len(values) > 1 {
			return values[1]
		}
	}

	var empty T
	return []T{empty}[0]
}
