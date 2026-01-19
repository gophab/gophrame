package feishu

import (
	"context"
	"time"

	"github.com/gophab/gophrame/core/consts"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security/server"
	"github.com/gophab/gophrame/core/social"
	"github.com/gophab/gophrame/core/social/feishu/config"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/errors"

	"fmt"
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/json"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkspeech_to_text "github.com/larksuite/oapi-sdk-go/v3/service/speech_to_text/v1"
	"github.com/patrickmn/go-cache"
)

type FeishuService struct {
	apps      sync.Map
	jsTickets *cache.Cache
}

func (s *FeishuService) createApp(appId string) *lark.Client {
	if appId == "" {
		appId = config.Setting.AppId
	}

	appSecret := ""
	if appId == config.Setting.AppId {
		appSecret = config.Setting.AppSecret
	} else if len(config.Setting.Apps) > 0 {
		for _, app := range config.Setting.Apps {
			if app.AppId == appId {
				appSecret = app.AppSecret
				break
			}
		}
	}

	if appSecret != "" {
		var client = lark.NewClient(appId, appSecret)
		s.apps.Store(appId, client)
		return client
	}
	return nil
}

func (s *FeishuService) GetApp(appId string) *lark.Client {
	if appId == "" {
		appId = config.Setting.AppId
	}

	if app, b := s.apps.Load(appId); b {
		return app.(*lark.Client)
	}

	return s.createApp(appId)
}

func (s *FeishuService) GetJsTicket(appId string) (string, bool) {
	if s.jsTickets == nil {
		s.jsTickets = cache.New(2*time.Hour, 5*time.Minute)
	}

	if ticket, b := s.jsTickets.Get(appId); b {
		return ticket.(string), ticket.(string) != "ERROR"
	}

	if app := s.GetApp(appId); app != nil {
		if apiResp, err := app.Post(context.Background(), "/open-apis/jssdk/ticket/get", nil, larkcore.AccessTokenTypeTenant); err == nil {
			resp := &GetTicketResp{ApiResp: apiResp}
			logger.Debug("Feishu POST result: ", apiResp.String())
			if err := json.Json(string(apiResp.RawBody), resp); err == nil {
				if resp.Data != nil {
					s.jsTickets.Set(appId, resp.Data.Ticket, time.Duration(resp.Data.ExpireIn)*time.Second)
					return resp.Data.Ticket, true
				}
			} else {
				logger.Error("Feishu get ticket error: ", err.Error())
				return "", false
			}
		} else {
			logger.Error("Feishu get ticket error: ", err.Error())
			s.jsTickets.SetDefault(appId, "ERROR")
			return "ERROR", false
		}
	}
	return "", false
}

type GetTicketRespData struct {
	ExpireIn int64  `json:"expire_in"`
	Ticket   string `json:"ticket"`
}

type GetTicketResp struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *GetTicketRespData `json:"data"` // 业务数据
}

func (s *FeishuService) GetSignature(ctx context.Context, appId string, url string, nonce string, timestamp int64) (*map[string]any, error) {
	if ticket, b := s.GetJsTicket(appId); b {
		if nonce == "" {
			nonce = object.QuickRandom(10)
		}
		if timestamp == 0 {
			timestamp = time.Now().Unix() * 1000
		}

		signature := util.SHA1(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, url))

		result := make(map[string]any)
		result["signature"] = signature
		result["appId"] = appId
		result["nonceStr"] = nonce
		result["timestamp"] = timestamp
		result["url"] = url

		return &result, nil
	}

	return nil, nil
}

func (s *FeishuService) SpeechToText(ctx context.Context, appId string, speech string) (string, error) {
	app := s.GetApp(appId)
	if app != nil {
		request := larkspeech_to_text.NewFileRecognizeSpeechReqBuilder().
			Body(larkspeech_to_text.NewFileRecognizeSpeechReqBodyBuilder().
				Speech(larkspeech_to_text.NewSpeechBuilder().
					Speech(speech).
					Build()).
				Config(larkspeech_to_text.NewFileConfigBuilder().
					FileId(object.QuickRandom(16)).
					Format("pcm").
					EngineType("16k_auto").
					Build()).
				Build()).
			Build()
		if resp, err := app.SpeechToText.Speech.FileRecognize(ctx, request); err != nil {
			return "", err
		} else {
			return *resp.Data.RecognitionText, nil
		}
	}

	return "", errors.New("No such appId")
}

func (s *FeishuService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *social.SocialUser {
	var appId string = ""
	segments := strings.Split(socialChannelId, ":")
	if len(segments) > 1 {
		appId = segments[1]
	}

	if appId == "" {
		value := ctx.Value(server.AppIdContextKey)
		if value != nil {
			segments = strings.Split(value.(string), ":")
			appId = segments[0]
		}
	}

	if app := s.GetApp(appId); app != nil {
		request := larkauthen.NewCreateAccessTokenReqBuilder().
			Body(larkauthen.NewCreateAccessTokenReqBodyBuilder().
				GrantType("authorization_code").
				Code(code).
				Build()).
			Build()
		if user, err := app.Authen.AccessToken.Create(ctx, request); err == nil {
			if user != nil {
				mobile := user.Data.Mobile
				if mobile != nil {
					// 截取mobile的最后11位数
					value := util.SubString(*mobile, len(*mobile)-11, 11)
					mobile = &value
				}
				var result = &social.SocialUser{
					Mobile:   mobile,
					Name:     user.Data.Name,
					Email:    user.Data.Email,
					Avatar:   user.Data.AvatarUrl,
					NickName: user.Data.EnName,
					OpenId:   user.Data.OpenId,
					Status:   util.IntAddr(consts.STATUS_VALID),
				}
				result.SetSocialId("fs", *user.Data.UnionId)
				return result
			}
		} else {
			logger.Warn("User from code error: ", code, err.Error())
		}
	}

	return nil
}
