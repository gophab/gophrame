package config

import (
	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"
)

type WebsocketSetting struct {
	Enabled               bool  `json:"enabled"`
	BufferSize            int   `json:"bufferSize" yaml:"bufferSize"`
	MaxMessageSize        int64 `json:"maxMessageSize" yaml:"maxMessageSize"`
	PingPeriod            int   `json:"pingPeriod" yaml:"pingPeriod"`
	HeartbeatFailMaxTimes int   `json:"heartbeatFailMaxTimes" yaml:"heartbeatFialMaxTimes"`
	ReadDeadline          int   `json:"readDeadline" yaml:"readDeadline"`
	WriteDeadline         int   `json:"writeDeadline" yaml:"writeDeadline"`
}

var Setting *WebsocketSetting = &WebsocketSetting{
	Enabled:               false,
	BufferSize:            20480,
	MaxMessageSize:        65535,
	PingPeriod:            20,
	HeartbeatFailMaxTimes: 4,
	ReadDeadline:          100,
	WriteDeadline:         35,
}

func init() {
	logger.Debug("Register Websocket Config")
	config.RegisterConfig("websocket", Setting, "Websocket Settings")
}
