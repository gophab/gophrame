package config

import (
	"github.com/gophab/gophrame/core/config"
	"github.com/gophab/gophrame/core/logger"
)

type AppSetting struct {
	CorpId        string `json:"corpId" yaml:"corpId"`
	AgentId       int    `json:"agentId" yaml:"agentId"`
	AppId         string `json:"appId" yaml:"appId"`
	AppSecret     string `json:"appSecret" yaml:"appSecret"`
	MessageToken  string `json:"messageToken" yaml:"messageToken"`
	MessageAESKey string `json:"messageAESKey" yaml:"messageAESKey"`
}

type DingtalkSetting struct {
	//AgentSetting
	CorpId        string       `json:"corpId" yaml:"corpId"`
	AgentId       int          `json:"agentId" yaml:"agentId"`
	AppId         string       `json:"appId" yaml:"appId"`
	AppSecret     string       `json:"appSecret" yaml:"appSecret"`
	MessageToken  string       `json:"messageToken" yaml:"messageToken"`
	MessageAESKey string       `json:"messageAESKey" yaml:"messageAESKey"`
	Enabled       bool         `json:"enabled" yaml:"enabled"`
	Apps          []AppSetting `json:"agents" yaml:"agents"`
}

var Setting *DingtalkSetting = &DingtalkSetting{
	Enabled: false,
}

func init() {
	logger.Debug("Register Social Config - Dingtalk")
	config.RegisterConfig("social.dingtalk", Setting, "Dingtalk Settings")
}
