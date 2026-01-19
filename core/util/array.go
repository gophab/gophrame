package util

import (
	"errors"
	"reflect"
)

func Contains(container any, obj any) (bool, error) {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if containerValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if containerValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in container")
}

func ContainsAny(container any, objs []any) (bool, error) {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if b, _ := Contains(objs, containerValue.Index(i).Interface()); b {
				return true, nil
			}
		}
	case reflect.Map:
		for _, obj := range objs {
			if containerValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
				return true, nil
			}
		}
	}
	return false, errors.New("not in container")
}

func ContainsAll(container any, objs []any) (bool, error) {
	for _, obj := range objs {
		if b, _ := Contains(container, obj); !b {
			return false, errors.New("not in container")
		}
	}
	return true, nil
}

func Filter(container any, objs []any) (any, error) {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		result := reflect.New(reflect.TypeOf(container))
		for i := 0; i < containerValue.Len(); i++ {
			if b, _ := Contains(objs, containerValue.Index(i).Interface()); !b {
				result = reflect.Append(result, containerValue.Index(i))
			}
		}
		return result, nil
	case reflect.Map:
		result := reflect.New(reflect.TypeOf(container))
		for _, key := range containerValue.MapKeys() {
			if b, _ := Contains(objs, key.Interface()); !b {
				value := containerValue.MapIndex(key)
				if !value.IsValid() {
					result.SetMapIndex(key, value)
				}
			}
		}
		return result, nil
	case reflect.Struct:
		result := reflect.New(reflect.TypeOf(container))
		containerType := reflect.TypeOf(container)
		for i := 0; i < containerType.NumField(); i++ {
			field := containerType.Field(i)
			if b, _ := Contains(objs, field.Name); !b {
				result.Field(i).Set(containerValue.Field(i))
			}
		}
		return result, nil
	}
	return nil, errors.New("unsupported container")
}

func First[T any](items []T, filter func(T) bool) T {
	for _, item := range items {
		if filter(item) {
			return item
		}
	}
	var empty T
	return empty
}
