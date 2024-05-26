package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type SnowFlakeSetting struct {
	MachineId int64 `json:"machineId" yaml:"machineId"`
}

var Setting *SnowFlakeSetting = &SnowFlakeSetting{
	MachineId: 13,
}

func init() {
	logger.Debug("Register Snowflake Config")
	config.RegisterConfig("snowflake", Setting, "Snowflake Settings")
}
