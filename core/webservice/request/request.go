package request

import (
	"errors"
	"io"
	"net/http/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Parameter string

func (f Parameter) Exist() bool {
	return string(f) != string([]byte{0x1E})
}

func (f Parameter) MustUint8() (uint8, error) {
	v, err := strconv.ParseUint(f.DefaultString(""), 10, 8)
	return uint8(v), err
}

func (f Parameter) MustInt() (int, error) {
	v, err := strconv.ParseInt(f.DefaultString(""), 10, 0)
	return int(v), err
}

func (f Parameter) MustInt64() (int64, error) {
	v, err := strconv.ParseInt(f.DefaultString(""), 10, 64)
	return int64(v), err
}

func (f Parameter) MustFloat64() (float64, error) {
	v, err := strconv.ParseFloat(f.DefaultString(""), 64)
	return float64(v), err
}

func (f Parameter) MustBool() (bool, error) {
	v, err := strconv.ParseBool(f.DefaultString(""))
	return v, err
}

func (f Parameter) MustString() (string, error) {
	if f.Exist() {
		return string(f), nil
	}
	return "", errors.New("参数为空")
}

func (f Parameter) DefaultUint8(defaultValue uint8) uint8 {
	if f.Exist() {
		if v, err := strconv.ParseUint(f.DefaultString("0"), 10, 8); err == nil {
			return uint8(v)
		}
	}
	return defaultValue
}

func (f Parameter) DefaultInt(defaultValue int) int {
	if f.Exist() {
		if v, err := strconv.ParseInt(f.DefaultString("0"), 10, 0); err == nil {
			return int(v)
		}
	}
	return defaultValue
}

func (f Parameter) DefaultInt64(defaultValue int64) int64 {
	if f.Exist() {
		if v, err := strconv.ParseInt(f.DefaultString("0"), 10, 64); err == nil {
			return v
		}
	}
	return defaultValue
}

func (f Parameter) DefaultFloat64(defaultValue float64) float64 {
	if f.Exist() {
		if v, err := strconv.ParseFloat(f.DefaultString("0"), 64); err == nil {
			return v
		}
	}
	return defaultValue
}

func (f Parameter) DefaultBool(defaultValue bool) bool {
	if f.Exist() {
		if v, err := strconv.ParseBool(f.DefaultString("0")); err == nil {
			return v
		}
	}
	return defaultValue
}

func (f Parameter) DefaultString(defaultValue string) string {
	if f.Exist() {
		return string(f)
	}
	return defaultValue
}

func (f Parameter) Uint8() uint8 {
	v, _ := f.MustUint8()
	return v
}

func (f Parameter) Int() int {
	v, _ := f.MustInt()
	return v
}

func (f Parameter) Int64() int64 {
	v, _ := f.MustInt64()
	return v
}

func (f Parameter) Float64() float64 {
	v, _ := f.MustFloat64()
	return v
}

func (f Parameter) Bool() bool {
	v, _ := f.MustBool()
	return v
}

func (f Parameter) String() *string {
	if f.Exist() {
		result := string(f)
		return &result
	}
	return nil
}

func Header(c *gin.Context, key string) Parameter {
	v := c.GetHeader(key)
	if v != "" {
		return Parameter(v)
	} else {
		return Parameter([]byte{0x1E})
	}
}

func Param(c *gin.Context, key string) Parameter {
	v, b := c.Params.Get(key)
	if !b {
		v, b = c.GetQuery(key)
		if !b {
			v, b = c.GetPostForm(key)
		}
	}

	if b {
		return Parameter(v)
	} else {
		return Parameter([]byte{0x1E})
	}
}

func ShouldParam(c *gin.Context, key string, p *Parameter) bool {
	v, b := c.Params.Get(key)
	if !b {
		v, b = c.GetQuery(key)
		if !b {
			v, b = c.GetPostForm(key)
		}
	}

	if b {
		*p = Parameter(v)
	}
	return b
}

func MustParam(c *gin.Context, key string) (Parameter, error) {
	v, b := c.Params.Get(key)
	if !b {
		v, b = c.GetQuery(key)
		if !b {
			v, b = c.GetPostForm(key)
		}
	}

	if b {
		return Parameter(v), nil
	} else {
		return Parameter([]byte{0x1E}), errors.New("no parameter '" + key + "'")
	}
}

type Parameters []Parameter

func Params(c *gin.Context, key string) Parameters {
	v, b := c.Params.Get(key)
	if !b {
		v, b = c.GetQuery(key)
		if !b {
			v, b = c.GetPostForm(key)
		}
	}

	if b {
		result := []Parameter{}
		for _, s := range v {
			result = append(result, Parameter(s))
		}
		return result
	} else {
		return nil
	}
}

func ShouldParams(c *gin.Context, key string, p *Parameters) bool {
	v, b := c.GetQueryArray(key)
	if !b {
		v, b = c.GetPostFormArray(key)
	}

	if b {
		for _, s := range v {
			*p = append(*p, Parameter(s))
		}
	}
	return b
}

func MustParams(c *gin.Context, key string) (Parameters, error) {
	v, b := c.GetQueryArray(key)
	if !b {
		v, b = c.GetPostFormArray(key)
	}

	if b {
		result := []Parameter{}
		for _, s := range v {
			result = append(result, Parameter(s))
		}
		return result, nil
	} else {
		return nil, errors.New("no parameter '" + key + "'")
	}
}

func (f Parameters) MustUint8() ([]uint8, error) {
	result := []uint8{}
	for _, v := range f {
		if ui, err := v.MustUint8(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) MustInt() ([]int, error) {
	result := []int{}
	for _, v := range f {
		if ui, err := v.MustInt(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) MustInt64() ([]int64, error) {
	result := []int64{}
	for _, v := range f {
		if ui, err := v.MustInt64(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) MustFloat64() ([]float64, error) {
	result := []float64{}
	for _, v := range f {
		if ui, err := v.MustFloat64(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) MustBool() ([]bool, error) {
	result := []bool{}
	for _, v := range f {
		if ui, err := v.MustBool(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) MustString() ([]string, error) {
	result := []string{}
	for _, v := range f {
		if ui, err := v.MustString(); err != nil {
			return nil, err
		} else {
			result = append(result, ui)
		}
	}

	return result, nil
}

func (f Parameters) DefaultUint8(defaultValue uint8) []uint8 {
	result := []uint8{}
	for _, v := range f {
		result = append(result, v.DefaultUint8(defaultValue))
	}

	return result
}

func (f Parameters) DefaultInt(defaultValue int) []int {
	result := []int{}
	for _, v := range f {
		result = append(result, v.DefaultInt(defaultValue))
	}

	return result
}

func (f Parameters) DefaultInt64(defaultValue int64) []int64 {
	result := []int64{}
	for _, v := range f {
		result = append(result, v.DefaultInt64(defaultValue))
	}

	return result
}

func (f Parameters) DefaultFloat64(defaultValue float64) []float64 {
	result := []float64{}
	for _, v := range f {
		result = append(result, v.DefaultFloat64(defaultValue))
	}

	return result
}

func (f Parameters) DefaultBool(defaultValue bool) []bool {
	result := []bool{}
	for _, v := range f {
		result = append(result, v.DefaultBool(defaultValue))
	}

	return result
}

func (f Parameters) Uint8() []uint8 {
	v, _ := f.MustUint8()
	return v
}

func (f Parameters) Int() []int {
	v, _ := f.MustInt()
	return v
}

func (f Parameters) Int64() []int64 {
	v, _ := f.MustInt64()
	return v
}

func (f Parameters) Float64() []float64 {
	v, _ := f.MustFloat64()
	return v
}

func (f Parameters) Bool() []bool {
	v, _ := f.MustBool()
	return v
}

func (f Parameters) String() []*string {
	result := []*string{}
	for _, v := range f {
		result = append(result, v.String())
	}

	return result
}

func Dump(writer io.Writer, context *gin.Context) error {
	data, err := httputil.DumpRequest(context.Request, true)
	if err != nil {
		return err
	}

	for header, values := range context.Request.Header {
		for _, value := range values {
			writer.Write([]byte("\n" + header + ": " + value + "\n"))
		}
	}

	writer.Write(data)
	return nil
}
