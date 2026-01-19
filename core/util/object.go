package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// 按名字和类型合并（&a,b）b->a: field值取不为空的
func MergeFields(dest any, src any, fields ...string) (err error) {
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
	if bt.Kind() == reflect.Ptr {
		bv = bv.Elem()
		bt = bv.Type()
	}

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
			if f.IsNil() || f.IsZero() || (!bValue.IsNil() && !bValue.IsZero()) {
				f.Set(bValue)
			}
		}
	}
	return
}

// 按名字和类型复制（&a,b）b->a: a.field为空的取b.field
func FillFields(dest any, src any, fields ...string) (err error) {
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
	if bt.Kind() == reflect.Ptr {
		bv = bv.Elem()
		bt = bv.Type()
	}

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
			if f.IsNil() || f.IsZero() {
				f.Set(bValue)
			}
		}
	}
	return
}

// 按名字和类型复制（&a,b）b->a: a.field为空的取b.field
func PatchFields(dest any, src any, fields ...string) (err error) {
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
	if bt.Kind() == reflect.Ptr {
		bv = bv.Elem()
		bt = bv.Type()
	}

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
			if !bValue.IsNil() && !bValue.IsZero() {
				f.Set(bValue)
			}
		}
	}
	return
}

// 按名字和类型复制（&a,b）b->a
func CopyFields(dest any, src any, fields ...string) (err error) {
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
	if bt.Kind() == reflect.Ptr {
		bv = bv.Elem()
		bt = bv.Type()
	}

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
func CopyFieldsExcept(dest any, src any, fields ...string) (err error) {
	at := reflect.TypeOf(dest)
	av := reflect.ValueOf(dest)
	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return
	}
	av = reflect.ValueOf(av.Interface())

	bt := reflect.TypeOf(src)
	bv := reflect.ValueOf(src)
	if bt.Kind() == reflect.Ptr {
		bv = bv.Elem()
		bt = bv.Type()
	}

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

func MapStruct(src map[string]any, dest any) {
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

func BoolAddr(b bool) *bool {
	return &b
}

func Struct2Map(src any) map[string]any {
	var result = make(map[string]any)
	if data, err := json.Marshal(src); err == nil {
		// Use json to
		_ = json.Unmarshal(data, &result)
	}
	return result
}

func GetRecordFieldValue[T any](record any, path string, v T) T {
	if field, ok := GetRecordField(record, path); ok {
		if data, err := json.Marshal(field); err == nil {
			// Use json to
			_ = json.Unmarshal(data, &v)
		}
	}
	return v
}

func GetRecordField(record any, path string) (any, bool) {
	if record == nil {
		return nil, false
	}

	segs := strings.SplitN(path, ".", 2)
	if len(segs) == 2 {
		if node, b := GetRecordField(record, segs[0]); b {
			return GetRecordField(node, segs[1])
		}
	} else {
		switch reflect.TypeOf(record).Kind() {
		case reflect.Map:
			value := reflect.ValueOf(record).MapIndex(reflect.ValueOf(path))
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return value.Interface(), true
			}
		case reflect.Array:
			if index, err := strconv.ParseInt(path, 10, 32); err == nil {
				array := reflect.ValueOf(record)
				if 0 <= index && index <= int64(array.Len()) {
					value := reflect.ValueOf(record).Index(int(index))
					if value.IsValid() && !value.IsNil() && !value.IsZero() {
						return value.Interface(), true
					}
				}
			}
		case reflect.Struct:
			value := reflect.ValueOf(record)
			field := value.FieldByName(path)
			if field.IsValid() && !field.IsNil() && !field.IsZero() {
				return field.Interface(), true
			}
		case reflect.Ptr, reflect.Interface:
			value := reflect.ValueOf(record).Elem()
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return GetRecordField(value.Interface(), path)
			}
		default:
		}
	}
	return nil, false
}

func MakeRecord(path string, value any) map[string]any {
	var result map[string]any
	var segs = strings.Split(path, ".")
	slices.Reverse(segs)
	for _, seg := range segs {
		result = make(map[string]any)
		result[seg] = value
		value = result
	}
	return result
}

func RemoveRecordField(value map[string]any, path string) {
	var segs = strings.Split(path, ".")
	slices.Reverse(segs)
	for i, seg := range segs {
		if i == len(segs)-1 {
			delete(value, seg)
		} else {
			if v, b := value[seg].(map[string]any); b {
				value = v
			} else {
				value = make(map[string]any)
			}
		}
	}
}
