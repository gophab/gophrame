package feishu

import (
	"context"
	"time"

	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security"
	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/social/feishu/config"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/errors"

	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/wjshen/gophrame/core/json"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/gin-gonic/gin"
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

func (s *FeishuService) GetSignature(ctx context.Context, appId string, url string, nonce string, timestamp int64) (*map[string]interface{}, error) {
	if ticket, b := s.GetJsTicket(appId); b {
		if nonce == "" {
			nonce = object.QuickRandom(10)
		}
		if timestamp == 0 {
			timestamp = time.Now().Unix() * 1000
		}

		signature := util.SHA1(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, url))

		result := make(map[string]interface{})
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

func (s *FeishuService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *domain.SocialUser {
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
				var result = &domain.SocialUser{
					Mobile:   mobile,
					Name:     user.Data.Name,
					Email:    user.Data.Email,
					Avatar:   user.Data.AvatarUrl,
					NickName: user.Data.EnName,
					OpenId:   user.Data.OpenId,
					Status:   &domain.STATUS_VALID,
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

type FeishuController struct {
	FeishuService *FeishuService
}

func (c *FeishuController) GetSignature(ctx *gin.Context) {
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
	if res, err := c.FeishuService.GetSignature(ctx, appId, path, nonce, timestamp); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, res)
	}
}

func (c *FeishuController) SpeechToText(ctx *gin.Context) {
	appId := request.Param(ctx, "appId").DefaultString(ctx.Request.Header.Get("X-App-Id"))
	if appId == "" {
		appId = config.Setting.AppId
	}

	//form表单
	data := ctx.Request.Form.Get("speech")
	if data == "" {
		response.FailMessage(ctx, 400, "上传文件失败")
		return
	}

	if res, err := c.FeishuService.SpeechToText(ctx, appId, data); err != nil {
		response.FailMessage(ctx, 400, err.Error())
		return
	} else {
		response.Success(ctx, map[string]string{"content": res})
	}
}

/**
 * 处理JSAPI的API路由
 */
func (e *FeishuController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		feishu := g.Group("/openapi/social/feishu")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			feishu.GET("/signature", e.GetSignature)                                  // 发送
			feishu.POST("/speech-text", security.HandleTokenVerify(), e.SpeechToText) // 发送
		}
	}
}
