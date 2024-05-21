package config

import "time"

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
