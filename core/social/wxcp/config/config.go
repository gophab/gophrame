package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type AgentSetting struct {
	CorpId        string `json:"corpId" yaml:"corpId"`
	AgentId       int    `json:"agentId" yaml:"agentId"`
	AppSecret     string `json:"appSecret" yaml:"appSecret"`
	MessageToken  string `json:"messageToken" yaml:"messageToken"`
	MessageAESKey string `json:"messageAESKey" yaml:"messageAESKey"`
}

type WxcpSetting struct {
	//AgentSetting
	CorpId        string         `json:"corpId" yaml:"corpId"`
	AgentId       int            `json:"agentId" yaml:"agentId"`
	AppSecret     string         `json:"appSecret" yaml:"appSecret"`
	MessageToken  string         `json:"messageToken" yaml:"messageToken"`
	MessageAESKey string         `json:"messageAESKey" yaml:"messageAESKey"`
	Enabled       bool           `json:"enabled" yaml:"enabled"`
	Agents        []AgentSetting `json:"agents" yaml:"agents"`
}

var Setting *WxcpSetting = &WxcpSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Social Config - Wxcp")
	config.RegisterConfig("social.wxcp", Setting, "Wxcp Settings")
}
