package wxcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gophab/gophrame/core/consts"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security/server"
	"github.com/gophab/gophrame/core/social"
	"github.com/gophab/gophrame/core/social/wxcp/config"
	"github.com/gophab/gophrame/core/social/wxcp/jssdk"
	"github.com/gophab/gophrame/core/util"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/providers"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

type Work struct {
	*work.Work
	JSSdk      *jssdk.Client
	AgentJSSdk *jssdk.Client
	apiToken   *string
	expireAt   time.Time
}

func (w *Work) GetAPIAccessToken() (string, error) {
	if w.apiToken != nil && w.expireAt.After(time.Now()) {
		return *w.apiToken, nil
	}

	result, err := w.OAuth.Provider.GetAPIAccessToken()
	if err != nil {
		return "", err
	}

	w.apiToken = &result
	w.expireAt = time.Now().Add(7200 * time.Second)
	return result, err
}

func (w *Work) UserFromCode(code string) (*providers.User, error) {
	token, err := w.GetAPIAccessToken()
	if err != nil {
		return nil, err
	}

	userInfo, err := w.OAuth.Provider.GetUser(token, code)
	if err != nil {
		return nil, err
	}

	if userDetail, err := w.OAuth.Provider.GetUserByID(userInfo.UserID); err == nil {
		// weCom.ResponseGetUserByID is detail response
		return providers.NewUser(&object.HashMap{
			"id":       userDetail.UserID,
			"name":     userDetail.Name,
			"avatar":   userDetail.Avatar,
			"nickname": userDetail.Alias,
			"mobile":   userDetail.Mobile,
			"email":    userDetail.Email,
			"title":    userDetail.Position,
			"openID":   userDetail.OpenID,
		}, w.OAuth.Provider), nil
	} else {
		return providers.NewUser(&object.HashMap{
			"id":     userInfo.UserID,
			"openid": userInfo.OpenID,
		}, w.OAuth.Provider), nil
	}
}

type WxcpService struct {
	apps sync.Map
}

func (s *WxcpService) createApp(agent *config.AgentSetting) *Work {
	corpId := agent.CorpId
	if corpId == "" {
		corpId = config.Setting.CorpId
	}
	if app, err := work.NewWork(&work.UserConfig{
		CorpID:  corpId, // 小程序、公众号或者企业微信的appid
		AgentID: agent.AgentId,
		Secret:  agent.AppSecret,
		Token:   agent.MessageToken,
		AESKey:  agent.MessageAESKey,

		OAuth: work.OAuth{
			Callback: "http://localhost",
			Scopes:   nil,
		},

		//Cache:     cache,
		HttpDebug: true,
		Debug:     true,
		//"sandbox": true,
	}); err == nil {
		client, _ := jssdk.RegisterProvider(app)
		client.TicketEndpoint = "cgi-bin/get_jsapi_ticket"
		client.TicketType = "jsapi"

		agentClient, _ := jssdk.RegisterProvider(app)
		agentClient.TicketEndpoint = "cgi-bin/ticket/get"
		agentClient.TicketType = "agent_config"

		result := &Work{
			Work:       app,
			JSSdk:      client,
			AgentJSSdk: agentClient,
		}
		s.apps.Store(fmt.Sprintf("%s:%d", corpId, agent.AgentId), result)
		return result
	} else {
		logger.Error("企业微信服务创建失败：", corpId, agent.AgentId, err.Error())
	}
	return nil
}

func (s *WxcpService) getAgentSetting(corpId string, agentId int) *config.AgentSetting {
	if corpId == config.Setting.CorpId && agentId == config.Setting.AgentId {
		return &config.AgentSetting{
			CorpId:        config.Setting.CorpId,
			AgentId:       config.Setting.AgentId,
			AppSecret:     config.Setting.AppSecret,
			MessageToken:  config.Setting.MessageToken,
			MessageAESKey: config.Setting.MessageAESKey,
		}
	}

	if len(config.Setting.Agents) > 0 {
		for i, agent := range config.Setting.Agents {
			if agent.AgentId == agentId && (agent.CorpId == corpId || (agent.CorpId == "" && config.Setting.CorpId == corpId)) {
				return &config.Setting.Agents[i]
			}
		}
	}
	return nil
}

func (s *WxcpService) GetApp(corpId string, agentId int) *Work {
	if corpId == "" {
		corpId = config.Setting.CorpId
	}

	if agentId == 0 {
		agentId = config.Setting.AgentId
	}

	if app, b := s.apps.Load(fmt.Sprintf("%s:%d", corpId, agentId)); b {
		return app.(*Work)
	}

	if agentSetting := s.getAgentSetting(corpId, agentId); agentSetting != nil {
		return s.createApp(agentSetting)
	}

	return nil
}

func (s *WxcpService) GetSignature(ctx context.Context, corpId string, agentId int, url string, nonce string, timestamp int64) (*object.HashMap, error) {
	if app := s.GetApp(corpId, agentId); app != nil {
		return app.JSSdk.ConfigSignature(ctx, url, nonce, timestamp)
	}

	return nil, nil
}

func (s *WxcpService) GetAgentSignature(ctx context.Context, corpId string, agentId int, url string, nonce string, timestamp int64) (*object.HashMap, error) {
	if app := s.GetApp(corpId, agentId); app != nil {
		return app.AgentJSSdk.ConfigSignature(ctx, url, nonce, timestamp)
	}

	return nil, nil
}

func (s *WxcpService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *social.SocialUser {
	var corpId string = ""
	var agentId int = 0
	segments := strings.Split(socialChannelId, ":")
	if len(segments) > 1 {
		corpId = segments[1]
	}
	if len(segments) > 2 {
		agentId, _ = strconv.Atoi(segments[2])
	}

	if corpId == "" {
		value := ctx.Value(server.AppIdContextKey)
		if value != nil {
			segments = strings.Split(value.(string), ":")
			corpId = segments[0]
			if len(segments) > 1 {
				agentId, _ = strconv.Atoi(segments[1])
			}
		}
	}

	if app := s.GetApp(corpId, agentId); app != nil {
		if user, err := app.UserFromCode(code); err == nil {
			if user != nil && user.GetID() != "" {
				var result *social.SocialUser

				mobile := user.GetMobile()
				if mobile != "" {
					// 截取mobile的最后11位数
					mobile = util.SubString(mobile, len(mobile)-11, 11)
				}

				result = &social.SocialUser{
					Mobile:   util.StringAddr(mobile),
					Name:     util.StringAddr(user.GetName()),
					Email:    util.StringAddr(user.GetEmail()),
					Avatar:   util.StringAddr(user.GetAvatar()),
					NickName: util.StringAddr(user.GetNickname()),
					OpenId:   util.StringAddr(user.GetOpenID()),
					Title:    util.StringAddr(user.GetString("title", "")),
					Status:   util.IntAddr(consts.STATUS_VALID),
				}
				result.SetSocialId("ww", user.GetID()+"@"+corpId)
				return result
			}
		} else {
			logger.Warn("User from code error: ", code, err.Error())
		}
	}

	return nil
}
