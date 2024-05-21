package error_record

import "github.com/wjshen/gophrame/core/logger"

// ErrorDeal 记录错误
func ErrorDeal(err error) error {
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}
