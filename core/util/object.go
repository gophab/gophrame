package util

import (
	"fmt"
	"reflect"
)

// 按名字和类型复制（&a,b）b->a
func CopyFields(dest interface{}, src interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(dest)
	av := reflect.ValueOf(dest)
	bt := reflect.TypeOf(src)
	bv := reflect.ValueOf(src)
	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return
	}
	av = reflect.ValueOf(av.Interface())
	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}
	if len(_fields) == 0 {
		return
	}
	// 复制
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)
		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		}
	}
	return
}

// 按名字和类型复制（&a,b）b->a
func CopyFieldsExcept(dest interface{}, src interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(dest)
	av := reflect.ValueOf(dest)
	bt := reflect.TypeOf(src)
	bv := reflect.ValueOf(src)
	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return
	}
	av = reflect.ValueOf(av.Interface())
	// 要复制哪些字段
	_fields := make([]string, 0)
	for i := 0; i < bv.NumField(); i++ {
		_fields = append(_fields, bt.Field(i).Name)
	}
	if len(fields) > 0 {
		_fields = DeleteElements(_fields, fields)
	}
	if len(_fields) == 0 {
		return
	}
	// 复制
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)
		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		}
	}
	return
}

func DeleteElement(src []string, elem string) []string {
	j := 0
	for _, v := range src {
		if v != elem {
			src[j] = v
			j++
		}
	}
	return src[:j]
}

func DeleteElements(src []string, elems []string) []string {
	var elemMap = make(map[string]int)
	for _, v := range elems {
		elemMap[v] = 1
	}
	j := 0
	for _, v := range src {
		if _, b := elemMap[v]; !b {
			src[j] = v
			j++
		}
	}
	return src[:j]
}

func MapStruct(src map[string]interface{}, dest interface{}) {
	v := reflect.ValueOf(dest).Elem()
	for key, val := range src {
		structField := v.FieldByName(key)
		if !structField.IsValid() {
			continue
		}
		structFieldType := structField.Type()
		valType := reflect.ValueOf(val)
		if structFieldType == valType.Type() {
			structField.Set(valType)
		}
	}
}
