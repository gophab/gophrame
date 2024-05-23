package wxmp

import (
	"context"
	"net/url"
	"strings"
	"sync"

	"github.com/wjshen/gophrame/core/consts"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxmp/config"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/gin-gonic/gin"
)

type WxmpService struct {
	apps sync.Map
}

func (s *WxmpService) createApp(appSetting *config.AppSetting) *officialAccount.OfficialAccount {
	if app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  appSetting.AppId, // 公众号appid
		Secret: appSetting.AppSecret,
		Token:  appSetting.MessageToken,
		AESKey: appSetting.MessageAESKey,

		ResponseType: "authorization_code",

		//Cache:     cache,
		HttpDebug: true,
		Debug:     false,
		//"sandbox": true,
	}); err == nil {
		s.apps.Store(appSetting.AppId, app)
		return app
	} else {
		logger.Error("微信公众号服务创建失败：", appSetting.AppId, err.Error())
	}
	return nil
}

func (s *WxmpService) getAppSetting(appId string) *config.AppSetting {
	if appId == config.Setting.AppId {
		return &config.AppSetting{
			AppId:         config.Setting.AppId,
			AppSecret:     config.Setting.AppSecret,
			MessageToken:  config.Setting.MessageToken,
			MessageAESKey: config.Setting.MessageAESKey,
		}
	}

	if len(config.Setting.Apps) > 0 {
		for i, app := range config.Setting.Apps {
			if app.AppId == appId {
				return &config.Setting.Apps[i]
			}
		}
	}
	return nil
}

func (s *WxmpService) GetApp(appId string) *officialAccount.OfficialAccount {
	if appId == "" {
		appId = config.Setting.AppId
	}

	if mp, b := s.apps.Load(appId); b {
		return mp.(*officialAccount.OfficialAccount)
	}

	if appSetting := s.getAppSetting(appId); appSetting != nil {
		return s.createApp(appSetting)
	}

	return nil
}

type WxToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
}

func (s *WxmpService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *social.SocialUser {
	var appId string
	segments := strings.Split(socialChannelId, ":")
	if len(segments) > 1 {
		appId = segments[1]
	}

	if appId == "" {
		value := ctx.Value(server.AppIdContextKey)
		if value != nil {
			appId = value.(string)
		}
	}

	if app := s.GetApp(appId); app != nil {
		if tokenResponse, err := app.OAuth.TokenFromCode(code); err == nil && tokenResponse != nil {
			var wxToken WxToken
			if err := object.HashMapToStructure(tokenResponse, &wxToken); err != nil {
				return nil
			}
			if strings.Contains(wxToken.Scope, "snsapi_userinfo") {
				if user, err := app.OAuth.UserFromToken(wxToken.AccessToken, wxToken.OpenId); err == nil && user != nil {
					mobile := user.GetMobile()
					if mobile != "" {
						// 截取mobile的最后11位数
						mobile = util.SubString(mobile, len(mobile)-11, 11)
					}

					result := &social.SocialUser{
						Mobile:   util.StringAddr(mobile),
						Name:     util.StringAddr(user.GetName()),
						Email:    util.StringAddr(user.GetEmail()),
						Avatar:   util.StringAddr(user.GetAvatar()),
						NickName: util.StringAddr(user.GetNickname()),
						OpenId:   util.StringAddr(user.GetOpenID() + "@" + appId),
						Status:   util.IntAddr(consts.STATUS_VALID),
					}
					// 判断是否存在UNIONID
					result.SetSocialId("wx", wxToken.UnionId)
					return result
				}
			} else {
				if user, err := app.User.Get(ctx, wxToken.OpenId, "zh_CN"); err == nil && user != nil {
					result := &social.SocialUser{
						OpenId: util.StringAddr(user.OpenID + "@" + appId),
						Remark: util.StringAddr(user.Remark),
						Status: util.IntAddr(consts.STATUS_VALID),
					}

					if user.UnionID != "" {
						result.SetSocialId("wx", user.UnionID)
					} else if result.OpenId != nil {
						result.SetSocialId("wx", *result.OpenId)
					}
					return result
				}
			}
			result := &social.SocialUser{
				OpenId: util.StringAddr(wxToken.OpenId + "@" + appId),
			}
			if wxToken.UnionId != "" {
				result.SetSocialId("wx", wxToken.UnionId)
			} else if result.OpenId != nil {
				result.SetSocialId("wx", *result.OpenId)
			}
			return result
		}

	}
	return nil
}

func (s *WxmpService) GetSignature(ctx context.Context, appId string, url string, nonce string, timestamp int64) (*object.HashMap, error) {
	if app := s.GetApp(appId); app != nil {
		return app.JSSDK.ConfigSignature(ctx, url, nonce, timestamp)
	}

	return nil, nil
}

type WxmpController struct {
	WxmpService *WxmpService
}

func (c *WxmpController) GetSignature(ctx *gin.Context) {
	appId := request.Param(ctx, "appId").DefaultString(ctx.Request.Header.Get("X-App-Id"))
	uri, err := request.Param(ctx, "url").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}
	nonce := request.Param(ctx, "nonceStr").DefaultString("")
	timestamp := request.Param(ctx, "timestamp").DefaultInt64(0)

	if appId == "" {
		appId = config.Setting.AppId
	}

	path, _ := url.QueryUnescape(uri)
	if res, err := c.WxmpService.GetSignature(ctx, appId, path, nonce, timestamp); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, res)
	}
}

/**
 * 处理WEB验证码的API路由
 */
func (e *WxmpController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		wxmp := g.Group("/openapi/social/wxmp")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			wxmp.GET("/signature", e.GetSignature) // 发送
		}
	}
}
