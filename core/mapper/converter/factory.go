package converter

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	createdConvertersMu sync.Mutex
	createdConverters   = make(map[convertType]*Converter)
)

func NewConverter(src, dst interface{}) (*Converter, error) {
	dstTyp := reflect.TypeOf(dst)
	srcTyp := reflect.TypeOf(src)

	if c := newConverter(srcTyp, dstTyp, nil, true); c == nil {
		return nil, fmt.Errorf("[MAPPER]can't convert source type %s to destination type %s", srcTyp, dstTyp)
	} else {
		return c, nil
	}
}

func NewConverterOption(src, dst interface{}, option *StructOption) (*Converter, error) {
	dstTyp := reflect.TypeOf(dst)
	srcTyp := reflect.TypeOf(src)

	if c := newConverter(srcTyp, dstTyp, option.convert(), true); c == nil {
		return nil, fmt.Errorf("can't convert source type %s to destination type %s", srcTyp, dstTyp)
	} else {
		return c, nil
	}
}

func newConverter(srcTyp, dstTyp reflect.Type, option *structOption, lock bool) *Converter {
	if lock {
		createdConvertersMu.Lock()
		defer createdConvertersMu.Unlock()
	}

	dstTyp = dereferencedType(dstTyp)
	srcTyp = dereferencedType(srcTyp)

	cTyp := &convertType{srcTyp, dstTyp, option}
	if dc, ok := createdConverters[*cTyp]; ok {
		return dc
	}

	var c converter
	if c = newBasicConverter(cTyp); c == nil {
		switch sk, dk := srcTyp.Kind(), dstTyp.Kind(); {

		case sk == reflect.Struct && dk == reflect.Struct:
			c = newStructConverter(cTyp)

		case sk == reflect.Slice && dk == reflect.Slice:
			c = newSliceConverter(cTyp)

		case sk == reflect.Map && dk == reflect.Map:
			c = newMapConverter(cTyp)
		}
	}
	if c != nil {
		dc := &Converter{cTyp, c}
		createdConverters[*cTyp] = dc
		return dc
	}

	return nil
}
