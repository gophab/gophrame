package config

import (
	"time"

	"github.com/wjshen/gophrame/core/config"
	"github.com/wjshen/gophrame/core/logger"
)

type CasbinSetting struct {
	Enabled                bool          `json:"enabled"`
	AutoLoadPolicyInterval time.Duration `json:"autoLoadPolicyInterval" yaml:"autoLoadPolicyInterval"`
	TablePrefix            string        `json:"tablePrefix" yaml:"tablePrefix"`
	TableName              string        `json:"tableName" yaml:"tableName"`
	ModelConfig            string        `json:"modelConfig" yaml:"modelConfig"`
}

var Setting *CasbinSetting = &CasbinSetting{
	Enabled:                false,
	AutoLoadPolicyInterval: time.Second * 5,
}

func init() {
	logger.Debug("Register Casbin Config")
	config.RegisterConfig("casbin", Setting, "Casbin Settings")
}
