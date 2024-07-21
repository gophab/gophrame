package errors

import "github.com/gophab/gophrame/core/logger"

type Error struct {
	Name    string `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var system_errors = make(map[string]Error)

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

func MakeError(name string) *Error {
	if err, b := system_errors[name]; b {
		return &err
	}
	logger.Warn("Not valid error name: ", name)
	return nil
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
