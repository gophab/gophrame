package mapper

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gophab/gophrame/core/mapper/converter"
)

type Render func(interface{}, interface{})

var (
	mutex      sync.Mutex
	converters = make(map[string]*converter.Converter)
	renders    = make(map[string]Render)
)

func getConverter(src, dst interface{}) (*converter.Converter, error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()

		if converter, err := converter.NewConverter(src, dst); err == nil {
			converters[key] = converter
			return converter, nil
		}
	}

	return nil, fmt.Errorf("[MAPPER] can't convert source type %s to destination type %s", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
}

func getRender(src, dst interface{}) Render {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if render, ok := renders[key]; ok {
		return render
	}
	return nil
}

func RegisterRender[S, D any](src S, dst D, render func(interface{}, interface{})) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	renders[key] = render
}

func Map(src, dst interface{}) error {
	converter, err := getConverter(src, dst)
	if err != nil {
		return err
	}

	if err = converter.Convert(src, dst); err != nil {
		return err
	}

	if render := getRender(src, dst); render != nil {
		// delegate mapper
		render(src, dst)
	}
	return nil
}

func MapAsWithError[T any](src interface{}, dst T) (result T, err error) {
	converter, err := getConverter(src, dst)
	if err != nil {
		return dst, err
	}

	if err = converter.Convert(src, dst); err != nil {
		return dst, err
	}

	if render := getRender(src, dst); render != nil {
		// delegate mapper
		render(src, dst)
	}
	return dst, nil
}

func MapAs[T any](src interface{}, dst T) (result T) {
	converter, err := getConverter(src, dst)
	if err != nil {
		return dst
	}

	converter.Convert(src, dst)

	if render := getRender(src, dst); render != nil {
		// delegate mapper
		render(src, dst)
	}

	return dst
}

func MapOption(src, dst interface{}, option *converter.StructOption) (err error) {
	converter, err := getConverter(src, dst)
	if err != nil {
		return
	}

	if err = converter.Convert(src, dst); err != nil {
		return
	}

	if render := getRender(src, dst); render != nil {
		// delegate mapper
		render(src, dst)
	}

	return
}

func MapRender(src, dst interface{}, render func(interface{}, interface{})) (err error) {
	converter, err := getConverter(src, dst)
	if err != nil {
		return err
	}

	if err = converter.Convert(src, dst); err != nil {
		return
	}

	if render != nil {
		render(src, dst)
	}

	return
}

func MapOptionRender(src, dst interface{}, option *converter.StructOption, render func(interface{}, interface{})) (err error) {
	converter, err := getConverter(src, dst)
	if err != nil {
		return err
	}

	if err = converter.Convert(dst, src); err != nil {
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

	converter, err := getConverter(src[0], dst[0])
	if err != nil {
		return err
	}

	render := getRender(src[0], dst[0])
	for index, s := range src {
		if err = converter.Convert(&s, &dst[index]); err != nil {
			return
		}

		if render != nil {
			// delegate mapper
			render(&s, &dst[index])
		}
	}

	return
}

func MapArrayAs[S, D any](src []S, dst D) []D {
	if len(src) == 0 {
		return []D{}
	}

	result, _ := MapArrayAsWithError(src, dst)

	return result
}

func MapArrayAsWithError[S, D any](src []S, dst D) (result []D, err error) {
	if len(src) == 0 {
		return
	}

	result = make([]D, 0)
	for _, s := range src {
		result = append(result, MapAs(s, dst))
	}
	return
}

func MapArrayOption[S, D any](src []S, dst []D, option *converter.StructOption) (err error) {
	if len(src) == 0 {
		return
	}

	converter, err := getConverter(src[0], dst[0])
	if err != nil {
		return err
	}

	render := getRender(src[0], dst[0])
	for index, s := range src {
		if err = converter.Convert(&s, &dst[index]); err != nil {
			return
		}

		if render != nil {
			// delegate mapper
			render(&s, &dst[index])
		}
	}

	return
}

func MapArrayRender[S, D any](src []S, dst []D, render func(interface{}, interface{})) (err error) {
	if len(src) == 0 {
		return
	}

	converter, err := getConverter(src[0], dst[0])
	if err != nil {
		return err
	}

	for index, s := range src {
		if err = converter.Convert(&s, &dst[index]); err != nil {
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

	converter, err := getConverter(src[0], dst[0])
	if err != nil {
		return err
	}

	for index, s := range src {
		if err = converter.Convert(&s, &dst[index]); err != nil {
			return
		}

		if render != nil {
			render(&s, &dst[index])
		}
	}

	return
}
