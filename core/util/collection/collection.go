package collection

import (
	"errors"
	"reflect"
	"strings"

	"golang.org/x/exp/constraints"
)

func Map[T any](container interface{}, f func(interface{}) T) []T {
	result := make([]T, 0)
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			result = append(result, f(containerValue.Index(i).Interface()))
		}
	case reflect.Map:
		iter := containerValue.MapRange()
		for iter.Next() {
			result = append(result, f(iter.Value().Interface()))
		}
	}
	return result
}

func MapToSet[T constraints.Ordered](container interface{}, f func(interface{}) T) Set[T] {
	result := make(Set[T])
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			result.Add(f(containerValue.Index(i).Interface()))
		}
	case reflect.Map:
		iter := containerValue.MapRange()
		for iter.Next() {
			result.Add(f(iter.Value().Interface()))
		}
	}
	return result
}

func Join[T any](container interface{}, f func(interface{}) string, delimeter string) string {
	segs := Map(container, f)
	return strings.Join(segs, delimeter)
}

func Contains(container interface{}, obj interface{}) (bool, error) {
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

func ContainsAny(container interface{}, objs []interface{}) (bool, error) {
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

func ContainsAll(container interface{}, objs []interface{}) (bool, error) {
	for _, obj := range objs {
		if b, _ := Contains(container, obj); !b {
			return false, errors.New("not in container")
		}
	}
	return true, nil
}

func Filter(container interface{}, objs []interface{}) (interface{}, error) {
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
