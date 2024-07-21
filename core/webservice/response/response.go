package response

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/errors"

	"github.com/gin-gonic/gin"
)

const (
	//业务处理发生错误
	BusinessOccurredErrorCode int    = -4000000
	BusinessOccurredErrorMsg  string = "业务执行错误,请联系管理员处理"
	//服务器代码发生错误
	ServerOccurredErrorCode int    = -5000000
	ServerOccurredErrorMsg  string = "服务器内部发生代码执行错误,请联系开发者排查错误日志"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) *gin.Context {
	return Response(g.C, httpCode, errCode, data)
}

func Response(c *gin.Context, httpCode, errCode int, data interface{}) *gin.Context {
	if errCode != 0 {
		c.JSON(httpCode, gin.H{
			"code":    errCode,
			"message": errors.GetErrorMessage(errCode),
		})
	} else {
		c.Header("Content-Type", "application/json; charset=utf-8")
		if data != nil {
			c.JSON(httpCode, data)
		} else {
			c.Status(httpCode)
		}
	}

	return c
}

// 将json字符窜以标准json格式返回（例如，从redis读取json、格式的字符串，返回给浏览器json格式）
func ResponseWithJson(c *gin.Context, httpCode int, jsonStr string) *gin.Context {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(httpCode, jsonStr)
	return c
}

func Error(c *gin.Context, httpCode int) {
	c.AbortWithStatus(httpCode)
}

func ErrorCode(c *gin.Context, httpCode int, errCode int) {
	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	c.AbortWithStatusJSON(httpCode, gin.H{
		"code":    errCode,
		"message": errors.GetErrorMessage(errCode),
	})
}

func ErrorMessage(c *gin.Context, httpCode int, dataCode int, msg string) {
	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	c.AbortWithStatusJSON(httpCode, gin.H{
		"code":    dataCode,
		"message": msg,
	})
}

// 语法糖函数封装
// 仅提交对象数据，不Abort()
func OK(c *gin.Context, data interface{}) *gin.Context {
	return Response(c, http.StatusOK, 0, data)
}

// 标准失败响应
func Bad(c *gin.Context, errCode int, errMsg string) {
	ErrorMessage(c, http.StatusBadRequest, errCode, errMsg)
}

func NotFound(c *gin.Context, msg string) {
	ErrorMessage(c, http.StatusNotFound, http.StatusNotFound, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	ErrorMessage(c, http.StatusUnauthorized, http.StatusUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	ErrorMessage(c, http.StatusForbidden, http.StatusForbidden, msg)
}

func Unprocessable(c *gin.Context, msg string) {
	ErrorMessage(c, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, msg)
}

func NotAllowed(c *gin.Context, msg string) {
	ErrorMessage(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, msg)
}

// 返回成功: OK() + Abort()
func Success(c *gin.Context, data interface{}) {
	OK(c, data)
	c.Abort()
}

func SuccessWithHeader(c *gin.Context, data interface{}, headers map[string]string) {
	for k, v := range headers {
		c.Header(k, v)
	}
	Success(c, data)
}

// 失败的业务逻辑
func Fail(c *gin.Context) {
	Bad(c, BusinessOccurredErrorCode, BusinessOccurredErrorMsg)
}

func FailError(c *gin.Context, err *errors.Error) {
	Bad(c, err.Code, err.Message)
}

func FailCode(c *gin.Context, errCode int) {
	Bad(c, errCode, errors.GetErrorMessage(errCode))
}

func FailMessage(c *gin.Context, errCode int, errMsg string) {
	Bad(c, errCode, errMsg)
}

func stack() string {
	buf := make([]byte, 1024*8)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return ("\n" + string(buf))
}

// 系统错误
func Exception(c *gin.Context, errCode int, errMessage string, msg string) {
	// internal log
	logger.Error("Code: ", errCode, " Message: ", errMessage, "\nError: ", msg, " \nStack: ", stack())

	if errMessage == "" {
		errMessage = msg
	}
	ErrorMessage(c, http.StatusInternalServerError, errCode, errMessage)
}

// 系统执行代码错误
func SystemError(c *gin.Context, err error) {
	Exception(c, ServerOccurredErrorCode, ServerOccurredErrorMsg, err.Error())
}

func SystemFail(c *gin.Context, err error) {
	Exception(c, ServerOccurredErrorCode, ServerOccurredErrorMsg, err.Error())
}

func SystemFailError(c *gin.Context, e *errors.Error, err error) {
	if e != nil {
		Exception(c, e.Code, e.Message, err.Error())
	} else {
		Exception(c, ServerOccurredErrorCode, err.Error(), err.Error())
	}
}

func SystemErrorCode(c *gin.Context, errCode int) {
	Exception(c, errCode, errors.GetErrorMessage(errCode), "")
}

func SystemErrorMessage(c *gin.Context, errCode int, msg string) {
	errMessage := errors.GetErrorMessage(errCode)
	if errMessage == "" {
		errMessage = msg
	}
	Exception(c, errCode, errMessage, msg)
}

func Page(context *gin.Context, total int64, list interface{}) {
	if total > 0 && list != nil {
		context.Header("X-Total-Count", strconv.FormatInt(total, 10))
		Success(context, list)
	} else {
		Success(context, []any{})
	}
}
