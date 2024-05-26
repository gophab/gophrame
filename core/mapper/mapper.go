package mapper

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gophab/gophrame/core/mapper/converter"
)

var (
	mutex      sync.Mutex
	converters = make(map[string]*converter.Converter)
)

func Map(src, dst interface{}) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverter(dst, src); err != nil {
			return
		}
	}
	if err = converters[key].Convert(dst, src); err != nil {
		return
	}
	return
}

func MapOption(src, dst interface{}, option *converter.StructOption) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverterOption(dst, src, option); err != nil {
			return
		}
	}
	if err = converters[key].Convert(dst, src); err != nil {
		return
	}
	return
}

func MapRender(src, dst interface{}, render func(interface{}, interface{})) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverter(dst, src); err != nil {
			return
		}
	}

	if err = converters[key].Convert(dst, src); err != nil {
		return
	}

	if render != nil {
		render(src, dst)
	}

	return
}

func MapOptionRender(src, dst interface{}, option *converter.StructOption, render func(interface{}, interface{})) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverterOption(dst, src, option); err != nil {
			return
		}
	}

	if err = converters[key].Convert(dst, src); err != nil {
		return
	}

	if render != nil {
		render(src, dst)
	}

	return
}

func MapArray[S, D any](src []S, dst []D) (err error) {
	if len(src) == 0 {
		return
	}

	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src[0]).String(), reflect.TypeOf(dst[0]).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverter(dst, src); err != nil {
			return
		}
	}

	for index, s := range src {
		if err = converters[key].Convert(&dst[index], &s); err != nil {
			return
		}
	}

	return
}

func MapArrayOption[S, D any](src []S, dst []D, option *converter.StructOption) (err error) {
	if len(src) == 0 {
		return
	}

	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src[0]).String(), reflect.TypeOf(dst[0]).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverterOption(dst, src, option); err != nil {
			return
		}
	}

	for index, s := range src {
		if err = converters[key].Convert(&dst[index], &s); err != nil {
			return
		}
	}

	return
}

func MapArrayRender[S, D any](src []S, dst []D, render func(interface{}, interface{})) (err error) {
	if len(src) == 0 {
		return
	}

	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src[0]).String(), reflect.TypeOf(dst[0]).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverter(dst, src); err != nil {
			return
		}
	}

	for index, s := range src {
		if err = converters[key].Convert(&dst[index], &s); err != nil {
			return
		}

		if render != nil {
			render(&s, &dst[index])
		}
	}

	return
}

func MapArrayOptionRender[S, D any](src []S, dst []D, option *converter.StructOption, render func(interface{}, interface{})) (err error) {
	if len(src) == 0 {
		return
	}

	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src[0]).String(), reflect.TypeOf(dst[0]).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = converter.NewConverterOption(dst, src, option); err != nil {
			return
		}
	}

	for index, s := range src {
		if err = converters[key].Convert(&dst[index], &s); err != nil {
			return
		}

		if render != nil {
			render(&s, &dst[index])
		}
	}

	return
}
