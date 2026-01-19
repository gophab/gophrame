package errors

import (
	"fmt"

	"github.com/gophab/gophrame/core/logger"
)

type Error struct {
	Name    string `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"err,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

var system_errors = make(map[string]Error)

func New(args ...any) error {
	code := 0
	message := ""
	var err error
	for _, arg := range args {
		switch v := arg.(type) {
		case int, int32, int64:
			code = v.(int)
		case string:
			message = v
		case error:
			err = v
		}
	}

	if message == "" && err != nil {
		message = err.Error()
	}

	if message == "" {
		message = GetErrorMessage(code)
	}

	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func RegisterError(name string, code int, message string) {
	system_errors[name] = Error{
		Code:    code,
		Message: message,
	}
}

func RegisterErrors(errs []Error) {
	for _, err := range errs {
		system_errors[err.Name] = Error{
			Name:    err.Name,
			Code:    err.Code,
			Message: err.Message,
		}
	}
}

func FromCode(code int) *Error {
	return &Error{
		Code:    code,
		Message: GetErrorMessage(code),
	}
}

func FromName(name string) *Error {
	return MakeError(name)
}

func MakeError(name string) *Error {
	if err, b := system_errors[name]; b {
		return &err
	}
	logger.Warn("Not valid error name: ", name)
	return &Error{
		Code:    500,
		Message: "Unknown internal error: " + name,
	}
}

func ErrorCode(name string) int {
	if err, b := system_errors[name]; b {
		return err.Code
	}
	logger.Warn("Not valid error name: ", name)
	return 0
}

func ErrorMessage(name string) string {
	if err, b := system_errors[name]; b {
		return err.Message
	}
	logger.Warn("Not valid error name: ", name)
	return ""
}
