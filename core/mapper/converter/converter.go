package converter

import (
	"fmt"
	"reflect"
	"unsafe"
)

type convertType struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
	option *structOption
}

// converter can handle converting among convertible basic types,
// and struct-struct, slice-slice, map-map converting too.
// type with nested pointer is supported.

// all methods in converter are thread-safe.
// we can define a global variable to hold a converter and use it in any goroutine.
type converter interface {
	convert(dPtr, sPtr unsafe.Pointer)
}

type Converter struct {
	*convertType
	converter
}

func (d *Converter) Convert(dst, src interface{}) error {
	if dst == nil || src == nil || reflect.ValueOf(dst).IsNil() || reflect.ValueOf(src).IsNil() {
		return nil
	}

	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		return fmt.Errorf("[MAPPER]destination should be a pointer. [actual:%v]", dv.Type())
	}

	if dv.Type() != d.dstTyp {
		return fmt.Errorf("[MAPPER]invalid destination type. [expected:%v] [actual:%v]", d.dstTyp, dv.Type())
	}

	sv := dereferencedValue(src)
	if !sv.CanAddr() {
		return fmt.Errorf("[MAPPER]source should be a pointer. [actual:%v]", sv.Type())
	}

	if sv.Type() != d.srcTyp {
		return fmt.Errorf("[MAPPER]invalid source type. [expected:%v] [actual:%v]", d.srcTyp, sv.Type())
	}

	d.converter.convert(unsafe.Pointer(dv.UnsafeAddr()), unsafe.Pointer(sv.UnsafeAddr()))
	return nil
}
