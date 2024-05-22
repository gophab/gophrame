package starter

import (
	"github.com/wjshen/gophrame/core/logger"
)

var starters = make([]func(), 0)

func RegisterStarter(f func()) {
	starters = append(starters, f)
}

func Start() {
	logger.Info("Starting starters...")
	for _, s := range starters {
		s()
	}
}
