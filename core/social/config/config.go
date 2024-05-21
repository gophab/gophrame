package config

import (
	DingtalkConfig "github.com/wjshen/gophrame/core/social/dingtalk/config"
	FeishuConfig "github.com/wjshen/gophrame/core/social/feishu/config"
	WxcpConfig "github.com/wjshen/gophrame/core/social/wxcp/config"
	WxmaConfig "github.com/wjshen/gophrame/core/social/wxma/config"
	WxmpConfig "github.com/wjshen/gophrame/core/social/wxmp/config"
)

type SocialSetting struct {
	Enabled  bool                            `json:"enabled" yaml:"enabled"`
	Wxmp     *WxmpConfig.WxmpSetting         `json:"wxmp" yaml:"wxmp"`
	Wxcp     *WxcpConfig.WxcpSetting         `json:"wxcp" yaml:"wxcp"`
	Wxma     *WxmaConfig.WxmaSetting         `json:"wxma" yaml:"wxma"`
	Feishu   *FeishuConfig.FeishuSetting     `json:"feishu" yaml:"feishu"`
	Dingtalk *DingtalkConfig.DingtalkSetting `json:"dingtalk" yaml:"dingtalk"`
}

var Setting *SocialSetting = &SocialSetting{
	Enabled:  false,
	Wxmp:     WxmpConfig.Setting,
	Wxcp:     WxcpConfig.Setting,
	Wxma:     WxmaConfig.Setting,
	Feishu:   FeishuConfig.Setting,
	Dingtalk: DingtalkConfig.Setting,
}
