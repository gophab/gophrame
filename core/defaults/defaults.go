package defaults

import (
	"github.com/creasty/defaults"

	"github.com/gophab/gophrame/core/logger"
)

func Default(v any) any {
	err := defaults.Set(v)
	if err != nil {
		logger.Warn("Object set defaults error: ", v, err.Error())
	}
	return v
}
