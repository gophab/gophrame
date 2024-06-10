package engine

import (
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/mock"
)

var (
	mutex  sync.Mutex
	engine *gin.Engine
)

func create() {
	mutex.Lock()
	if engine == nil {
		engine = gin.New()
		engine.Use(gin.Logger()) // 日志
		engine.Use(gin.Recovery())
	}
	mutex.Unlock()
}

func Init(debug bool) *gin.Engine {
	engine := Get()

	// 调试模式
	if debug {
		gin.SetMode(gin.DebugMode)
		pprof.Register(engine)

		engine.Use(mock.Mock())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return engine
}

func Get() *gin.Engine {
	if engine == nil {
		create()
	}

	return engine
}
