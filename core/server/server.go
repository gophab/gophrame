package server

import (
	"fmt"
	"net/http"

	"github.com/gophab/gophrame/core/engine"
	"github.com/gophab/gophrame/core/logger"

	"github.com/gophab/gophrame/core/server/config"
)

func Start() {
	// 服务器配置
	if config.Setting.Enabled {
		readTimeout := config.Setting.ReadTimeout
		writeTimeout := config.Setting.WriteTimeout

		endPoint := fmt.Sprintf("%s:%d", config.Setting.BindAddr, config.Setting.Port)
		maxHeaderBytes := 1 << 20

		server := &http.Server{
			Addr:           endPoint,
			Handler:        engine.Get(),
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		}

		logger.Info("Start http server listening: ", endPoint)

		_ = server.ListenAndServe()
	}
}
