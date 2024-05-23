package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/wjshen/gophrame/core/webservice/response"

	"github.com/gin-gonic/gin"
)

const (
	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "参数校验失败"

	// 表单参数验证器未通过时的错误
	ErrorsValidatorNotExists      string = "不存在的验证器"
	ErrorsValidatorTransInitFail  string = "validator的翻译器初始化错误"
	ErrorNotAllParamsIsBlank      string = "该接口不允许所有参数都为空,请按照接口要求提交必填参数"
	ErrorsValidatorBindParamsFail string = "验证器绑定参数失败"
)

// 参数校验错误
func ErrorParam(c *gin.Context, wrongParam interface{}) {
	response.FailMessage(c, ValidatorParamsCheckFailCode, ValidatorParamsCheckFailMsg)
	c.Abort()
}

// ValidatorError 翻译表单参数验证器出现的校验错误
func ValidatorError(c *gin.Context, err error) {
	if _, ok := err.(validator.ValidationErrors); ok {
		//wrongParam := RemoveTopStruct(errs.Translate(Trans))
		response.FailMessage(c, ValidatorParamsCheckFailCode, ValidatorParamsCheckFailMsg)
	} else {
		errStr := err.Error()
		// multipart:nextpart:eof 错误表示验证器需要一些参数，但是调用者没有提交任何参数
		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			response.FailMessage(c, ValidatorParamsCheckFailCode, ErrorNotAllParamsIsBlank)
		} else {
			response.FailMessage(c, ValidatorParamsCheckFailCode, ValidatorParamsCheckFailMsg)
		}
	}
	c.Abort()
}

// Trans 定义一个全局翻译器T
var Trans ut.Translator

// InitTrans 初始化表单参数验证器的翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		//初始化翻译器
		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		// locale 通常取决于 http 请求头的 'Accept-Language'
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		//注册翻译器
		//默认注册英文，en 注册英文 zh 注册中文
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

// RemoveTopStruct 将返回的结构体名去除掉，只留下需要的字段名
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.LastIndex(field, ".")+1:]] = err
	}
	return res
}
