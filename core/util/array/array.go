package array

import (
	"errors"
	"strings"
)

func Map[S, T any](list []S, f func(item S) T) []T {
	segs := make([]T, 0)
	for _, item := range list {
		segs = append(segs, f(item))
	}
	return segs
}

func Join[T any](list []T, f func(item T) string, delimeter string) string {
	segs := Map(list, f)
	return strings.Join(segs, delimeter)
}

func Filter[T any](list []T, f func(T) bool) ([]T, error) {
	result := make([]T, 0)
	for _, v := range list {
		if b := f(v); !b {
			result = append(result, v)
		}
	}
	return result, nil
}

func Contains[T comparable](list []T, obj T) (bool, error) {
	for _, v := range list {
		if v == obj {
			return true, nil
		}
	}

	return false, errors.New("not in container")
}

func ContainsAny[T comparable](list []T, objs []T) (bool, error) {
	for _, v := range list {
		if b, _ := Contains(objs, v); b {
			return true, nil
		}
	}

	return false, errors.New("not in container")
}

func ContainsAll[T comparable](list []T, objs []T) (bool, error) {
	for _, obj := range objs {
		if b, _ := Contains(list, obj); !b {
			return false, errors.New("not in container")
		}
	}
	return true, nil
}

func FilterAll[T comparable](list []T, objs []T) ([]T, error) {
	return Filter(list, func(v T) bool {
		b, _ := Contains(objs, v)
		return b
	})
}
