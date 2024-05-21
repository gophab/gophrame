package wxcp

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/social/wxcp/config"
	"github.com/wjshen/gophrame/core/social/wxcp/jssdk"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/errors"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/providers"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
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

func (s *WxcpService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *domain.SocialUser {
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
				var result *domain.SocialUser

				mobile := user.GetMobile()
				if mobile != "" {
					// 截取mobile的最后11位数
					mobile = util.SubString(mobile, len(mobile)-11, 11)
				}

				result = &domain.SocialUser{
					Mobile:   util.StringAddr(mobile),
					Name:     util.StringAddr(user.GetName()),
					Email:    util.StringAddr(user.GetEmail()),
					Avatar:   util.StringAddr(user.GetAvatar()),
					NickName: util.StringAddr(user.GetNickname()),
					OpenId:   util.StringAddr(user.GetOpenID()),
					Title:    util.StringAddr(user.GetString("title", "")),
					Status:   &domain.STATUS_VALID,
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

type WxcpController struct {
	WxcpService *WxcpService
}

func (c *WxcpController) GetSignature(ctx *gin.Context) {
	corpId := request.Param(ctx, "corpId").DefaultString(ctx.Request.Header.Get("X-Corp-Id"))
	agentId := request.Param(ctx, "agentId").DefaultString(ctx.Request.Header.Get("X-Agent-Id"))
	uri, err := request.Param(ctx, "url").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}
	nonce := request.Param(ctx, "nonceStr").DefaultString("")
	timestamp := request.Param(ctx, "timestamp").DefaultInt64(0)

	if corpId == "" {
		corpId = config.Setting.CorpId
	}

	agentIdInt := 0
	if agentId == "" {
		agentIdInt = config.Setting.AgentId
	} else {
		agentIdInt, _ = strconv.Atoi(agentId)
	}

	path, _ := url.QueryUnescape(uri)
	if res, err := c.WxcpService.GetSignature(ctx, corpId, agentIdInt, path, nonce, timestamp); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, res)
	}
}

func (c *WxcpController) GetAgentSignature(ctx *gin.Context) {
	corpId := request.Param(ctx, "corpId").DefaultString(ctx.Request.Header.Get("X-Corp-Id"))
	agentId := request.Param(ctx, "agentId").DefaultString(ctx.Request.Header.Get("X-Agent-Id"))
	uri, err := request.Param(ctx, "url").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}
	nonce := request.Param(ctx, "nonceStr").DefaultString("")
	timestamp := request.Param(ctx, "timestamp").DefaultInt64(0)

	if corpId == "" {
		corpId = config.Setting.CorpId
	}

	agentIdInt := 0
	if agentId == "" {
		agentIdInt = config.Setting.AgentId
	} else {
		agentIdInt, _ = strconv.Atoi(agentId)
	}

	path, _ := url.QueryUnescape(uri)
	if res, err := c.WxcpService.GetAgentSignature(ctx, corpId, agentIdInt, path, nonce, timestamp); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, res)
	}
}

/**
 * 处理WEB验证码的API路由
 */
func (e *WxcpController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		wxcp := g.Group("/openapi/social/wxcp")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			wxcp.GET("/signature", e.GetSignature)            // 发送
			wxcp.GET("/agent-signature", e.GetAgentSignature) // 发送
		}
	}
}
