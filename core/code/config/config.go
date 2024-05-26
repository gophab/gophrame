package config

import "time"

type RedisCodeStoreSetting struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Database  int    `json:"database" yaml:"database"`
	KeyPrefix string `json:"keyPrefix" yaml:"keyPrefix"`
}

type CacheCodeStoreSetting struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

type CodeStoreSetting struct {
	Enabled         bool                   `json:"enabled" yaml:"enabled"`
	RequestInterval time.Duration          `json:"requestInterval" yaml:"requestInterval"`
	ExpireIn        time.Duration          `json:"expireIn" yaml:"expireIn"`
	Cache           *CacheCodeStoreSetting `json:"cache"`
	Redis           *RedisCodeStoreSetting `json:"store"`
}
