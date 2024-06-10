package converter

import (
	"reflect"
	"unsafe"
)

type sliceConverter struct {
	*convertType
	*elemConverter
	sElemSize uintptr
	dElemSize uintptr
}

func newSliceConverter(convertType *convertType) (s converter) {
	c := &sliceConverter{
		convertType: convertType,
		sElemSize:   convertType.srcTyp.Elem().Size(),
		dElemSize:   convertType.dstTyp.Elem().Size(),
	}
	if convertType.srcTyp == convertType.dstTyp {
		s = c
	} else if elemConverter, ok := newElemConverter(convertType.srcTyp.Elem(), convertType.dstTyp.Elem(), convertType.option); ok {
		c.elemConverter = elemConverter
		s = c
	}
	return
}

// convert will overwrite the whole destination slice.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by Converter.Convert() and elemConverter.convert()
func (s *sliceConverter) convert(sPtr, dPtr unsafe.Pointer) {
	dSlice := (*sliceHeader)(dPtr)
	sSlice := (*sliceHeader)(sPtr)

	length := sSlice.Len
	dSlice.Len = length

	if dSlice.Cap < length || dSlice.Data == nil {
		newVal := reflect.MakeSlice(s.dstTyp, 0, length)
		dSlice.Data = unsafe.Pointer(newVal.Pointer())
		dSlice.Cap = length
	}

	if s.srcTyp == s.dstTyp {
		Copy(sSlice.Data, dSlice.Data, uintptr(length)*s.sElemSize)
		return
	}

	for dOffset, sOffset, i := uintptr(0), uintptr(0), 0; i < length; i++ {
		dElemPtr := unsafe.Pointer(uintptr(dSlice.Data) + dOffset)
		sElemPtr := unsafe.Pointer(uintptr(sSlice.Data) + sOffset)
		s.elemConverter.convert(sElemPtr, dElemPtr)
		dOffset += s.dElemSize
		sOffset += s.sElemSize
	}
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
