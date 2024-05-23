package wxma

import (
	"context"
	"strings"
	"sync"

	"github.com/wjshen/gophrame/core/consts"
	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/wxma/config"
	"github.com/wjshen/gophrame/core/util"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

type WxmaService struct {
	apps sync.Map
}

func (s *WxmaService) createApp(app *config.AppSetting) *miniProgram.MiniProgram {
	if mp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:  app.AppId, // 小程序appid
		Secret: app.AppSecret,
		Token:  app.MessageToken,
		AESKey: app.MessageAESKey,

		//Cache:     cache,
		HttpDebug: true,
		Debug:     false,
		//"sandbox": true,
	}); err == nil {
		s.apps.Store(app.AppId, mp)
		return mp
	}
	return nil
}

func (s *WxmaService) getAppSetting(appId string) *config.AppSetting {
	if appId == "" || appId == config.Setting.AppId {
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

func (s *WxmaService) GetApp(appId string) *miniProgram.MiniProgram {
	if appId == "" {
		appId = config.Setting.AppId
	}

	if ma, b := s.apps.Load(appId); b {
		return ma.(*miniProgram.MiniProgram)
	}

	if appSetting := s.getAppSetting(appId); appSetting != nil {
		return s.createApp(appSetting)
	}
	return nil
}

func (s *WxmaService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *social.SocialUser {
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
		if user, err := app.Auth.Session(ctx, code); err == nil && user != nil {
			result := social.SocialUser{
				OpenId: util.StringAddr(user.OpenID),
				Status: util.IntAddr(consts.STATUS_VALID),
			}
			if user.UnionID != "" {
				result.SetSocialId("wx", user.UnionID)
			} else {
				result.SetSocialId("wx", user.OpenID)
			}
			return &result
		}
	}

	return nil
}
