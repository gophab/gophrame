package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/wjshen/gophrame/core/consts"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/security/server"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/social/dingtalk/config"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/errors"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type Client struct {
	*http.Client
	Base      string
	AppKey    string
	AppSecret string
}

func (c *Client) Clone() *Client {
	return &Client{
		Client:    &http.Client{},
		Base:      c.Base,
		AppKey:    c.AppKey,
		AppSecret: c.AppSecret,
	}
}

type OpenAPIResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Unmarshallable interface {
	checkError() error
}

func (oar *OpenAPIResponse) checkError() error {
	var err error
	if oar.ErrCode != 0 {
		err = fmt.Errorf("errcode: %d\nerrmsg: %s", oar.ErrCode, oar.ErrMsg)
	}
	return err
}

func (c *Client) doRequest(request *http.Request, data interface{}) error {
	logger.Debug("Http Request ", request.Method, request.URL)
	if resp, err := c.Do(request); err == nil {
		if resp != nil {
			if resp.StatusCode >= 300 {
				logger.Error("Http Request error: ", resp.Status)
				return errors.New("Server Error: " + resp.Status)
			}

			defer resp.Body.Close()

			if content, err := ioutil.ReadAll(resp.Body); err == nil {
				logger.Debug("Http Request result: ", string(content))
				if err := json.Unmarshal(content, data); err != nil {
					return err
				}
				if v, b := data.(Unmarshallable); b {
					return v.checkError()
				}
			}
		}
		return nil
	} else {
		return err
	}
}

type GetAccessTokenRequest struct {
	// 应用id
	AppKey *string `json:"appKey,omitempty" xml:"appKey,omitempty"`
	// 应用密码
	AppSecret *string `json:"appSecret,omitempty" xml:"appSecret,omitempty"`
}

func (s GetAccessTokenRequest) String() string {
	return tea.Prettify(s)
}

func (s GetAccessTokenRequest) GoString() string {
	return s.String()
}

func (s *GetAccessTokenRequest) SetAppKey(v string) *GetAccessTokenRequest {
	s.AppKey = &v
	return s
}

func (s *GetAccessTokenRequest) SetAppSecret(v string) *GetAccessTokenRequest {
	s.AppSecret = &v
	return s
}

type GetAccessTokenResponse struct {
	*OpenAPIResponse
	// accessToken
	AccessToken string `json:"access_token,omitempty" xml:"accessToken,omitempty"`
	// 超时时间
	ExpiresIn int64 `json:"expires_in,omitempty" xml:"expireIn,omitempty"`
}

func (s GetAccessTokenResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessTokenResponse) GoString() string {
	return s.String()
}

func (s *GetAccessTokenResponse) SetAccessToken(v string) *GetAccessTokenResponse {
	s.AccessToken = v
	return s
}

func (s *GetAccessTokenResponse) SetExpiresIn(v int64) *GetAccessTokenResponse {
	s.ExpiresIn = v
	return s
}

func (c *Client) GetAccessToken(request *GetAccessTokenRequest) (*GetAccessTokenResponse, error) {
	query := url.Values{}
	query.Set("appkey", c.AppKey)
	query.Set("appsecret", c.AppSecret)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?%s", c.Base, "/gettoken", query.Encode()), nil)
	if err != nil {
		return nil, err
	}

	response := &GetAccessTokenResponse{}
	err = c.doRequest(req, response)
	return response, err
}

type GetJsTicketResponse struct {
	*OpenAPIResponse
	Ticket    string `json:"ticket,omitempty" xml:"ticket,omitempty"`        // accessToken
	ExpiresIn int64  `json:"expires_in,omitempty" xml:"expiresIn,omitempty"` // 超时时间
}

func (s GetJsTicketResponse) String() string {
	return tea.Prettify(s)
}

func (s GetJsTicketResponse) GoString() string {
	return s.String()
}

func (s *GetJsTicketResponse) SetTicket(v string) *GetJsTicketResponse {
	s.Ticket = v
	return s
}

func (s *GetJsTicketResponse) SetExpireIn(v int64) *GetJsTicketResponse {
	s.ExpiresIn = v
	return s
}

func (c *Client) GetJsTicket(accessToken string) (*GetJsTicketResponse, error) {
	query := url.Values{}
	query.Set("access_token", accessToken)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?%s", c.Base, "/get_jsapi_ticket", query.Encode()), nil)
	if err != nil {
		return nil, err
	}

	response := &GetJsTicketResponse{}
	err = c.doRequest(req, response)
	return response, err
}

type GetUserInfoResponseBody struct {
	UserId   string `json:"userid" xml:"userid"`
	DeviceId string `json:"device_id" xml:"deviceId"`
	Sys      bool   `json:"sys" xml:"sys"`
	SysLevel int    `json:"sys_level" xml:"sysLevel"`
	UnionId  string `json:"unionid" xml:"unionid"`
	Name     string `json:"name" xml:"name"`
}

func (s GetUserInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUserInfoResponseBody) GoString() string {
	return s.String()
}

type GetUserInfoResponse struct {
	*OpenAPIResponse
	Headers map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Result  *GetUserInfoResponseBody `json:"result,omitempty" xml:"result,omitempty" require:"true"`
}

func (c *Client) GetUserInfo(code string, accessToken string) (*GetUserInfoResponse, error) {
	query := url.Values{}
	query.Set("access_token", accessToken)

	if body, err := json.Marshal(map[string]string{"code": code}); err == nil {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s?%s", c.Base, "/topapi/v2/user/getuserinfo", query.Encode()), bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		response := &GetUserInfoResponse{}
		err = c.doRequest(req, response)
		return response, err
	} else {
		return nil, err
	}
}

type GetUserDetailResponseBody struct {
	UserId   string `json:"userid" xml:"userid"`
	UnionId  string `json:"unionid" xml:"unionid"`
	Name     string `json:"name" xml:"name"`
	Avatar   string `json:"avatar" xml:"avatar"`
	Mobile   string `json:"mobile" xml:"mobile"`
	Email    string `json:"email" xml:"email"`
	Remark   string `json:"remark" xml:"remark"`
	Title    string `json:"title" xml:"title"`
	DeviceId string `json:"device_id" xml:"deviceId"`
	Sys      bool   `json:"sys" xml:"sys"`
	SysLevel int    `json:"sys_level" xml:"sysLevel"`
}

func (s GetUserDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUserDetailResponseBody) GoString() string {
	return s.String()
}

type GetUserDetailResponse struct {
	*OpenAPIResponse
	Headers map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Result  *GetUserDetailResponseBody `json:"result,omitempty" xml:"result,omitempty" require:"true"`
}

func (c *Client) GetUserDetail(userId string, accessToken string) (*GetUserDetailResponse, error) {
	query := url.Values{}
	query.Set("access_token", accessToken)

	if body, err := json.Marshal(map[string]string{"userid": userId}); err == nil {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s?%s", c.Base, "/topapi/v2/user/get", query.Encode()), bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		response := &GetUserDetailResponse{}
		err = c.doRequest(req, response)
		return response, err
	} else {
		return nil, err
	}

}

type DingtalkService struct {
	apps      sync.Map
	tokens    *cache.Cache
	jsTickets *cache.Cache
}

func (s *DingtalkService) createApp(cfg *config.AppSetting) *Client {
	result := &Client{
		Client:    &http.Client{},
		Base:      "https://oapi.dingtalk.com",
		AppKey:    cfg.AppId,
		AppSecret: cfg.AppSecret,
	}
	s.apps.Store(cfg.AppId, result)
	return result
}

func (s *DingtalkService) getAppSetting(corpId string, agentId int) *config.AppSetting {
	if corpId == config.Setting.CorpId && agentId == config.Setting.AgentId {
		return &config.AppSetting{
			CorpId:        config.Setting.CorpId,
			AgentId:       config.Setting.AgentId,
			AppId:         config.Setting.AppId,
			AppSecret:     config.Setting.AppSecret,
			MessageToken:  config.Setting.MessageToken,
			MessageAESKey: config.Setting.MessageAESKey,
		}
	}

	if len(config.Setting.Apps) > 0 {
		for i, app := range config.Setting.Apps {
			if app.AgentId == agentId && (app.CorpId == corpId || (app.CorpId == "" && config.Setting.CorpId == corpId)) {
				return &config.Setting.Apps[i]
			}
		}
	}
	return nil
}

func (s *DingtalkService) getAppSettingByAppId(appId string) *config.AppSetting {
	if appId == config.Setting.AppId {
		return &config.AppSetting{
			CorpId:        config.Setting.CorpId,
			AgentId:       config.Setting.AgentId,
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

func (s *DingtalkService) GetApp(corpId string, agentId int) *Client {
	if corpId == "" {
		corpId = config.Setting.CorpId
	}

	if agentId == 0 {
		agentId = config.Setting.AgentId
	}

	if app, b := s.apps.Load(fmt.Sprintf("%s:%d", corpId, agentId)); b {
		return app.(*Client)
	}

	if appSetting := s.getAppSetting(corpId, agentId); appSetting != nil {
		return s.createApp(appSetting)
	}

	return nil
}

func (s *DingtalkService) GetAppByAppId(appId string) *Client {
	if appId == "" {
		appId = config.Setting.AppId
	}

	if app, b := s.apps.Load(appId); b {
		return app.(*Client)
	}

	if appSetting := s.getAppSettingByAppId(appId); appSetting != nil {
		return s.createApp(appSetting)
	}

	return nil
}

func (s *DingtalkService) GetAccessToken(appId string) (string, bool) {
	if s.tokens == nil {
		s.tokens = cache.New(2*time.Hour, 5*time.Minute)
	}

	if token, b := s.tokens.Get(appId); b {
		return token.(string), token.(string) != "ERROR"
	}

	if app := s.GetAppByAppId(appId); app != nil {
		request := &GetAccessTokenRequest{
			AppKey:    &app.AppKey,
			AppSecret: &app.AppSecret,
		}

		if resp, err := app.Clone().GetAccessToken(request); err == nil {
			s.tokens.Set(appId, resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second)
			return resp.AccessToken, true
		} else {
			logger.Error("Dingtalk get access error: ", err.Error())
			s.tokens.SetDefault(appId, "ERROR")
			return "ERROR", false
		}
	}
	return "", false
}

func (s *DingtalkService) GetJsTicket(appId string) (string, bool) {
	if s.jsTickets == nil {
		s.jsTickets = cache.New(2*time.Hour, 5*time.Minute)
	}

	if ticket, b := s.jsTickets.Get(appId); b {
		return ticket.(string), ticket.(string) != "ERROR"
	}

	if app := s.GetAppByAppId(appId); app != nil {
		if token, b := s.GetAccessToken(appId); b {
			if resp, err := app.Clone().GetJsTicket(token); err == nil {
				s.jsTickets.Set(appId, resp.Ticket, time.Duration(resp.ExpiresIn)*time.Second)
				return resp.Ticket, true
			} else {
				logger.Error("Dingtalk get ticket error: ", err.Error())
				s.jsTickets.SetDefault(appId, "ERROR")
				return "ERROR", false
			}
		}
	}
	return "", false
}

func (s *DingtalkService) GetSignature(ctx context.Context, appId string, url string, nonce string, timestamp int64) (*map[string]interface{}, error) {
	cfg := s.getAppSettingByAppId(appId)
	if cfg == nil {
		return nil, errors.New("应用不存在")
	}

	if ticket, b := s.GetJsTicket(appId); b {
		if nonce == "" {
			nonce = object.QuickRandom(10)
		}
		if timestamp == 0 {
			timestamp = time.Now().Unix()
		}

		signature := util.SHA1(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, url))

		result := make(map[string]interface{})
		result["signature"] = signature
		result["appId"] = appId
		result["corpId"] = cfg.CorpId
		result["agentId"] = cfg.AgentId
		result["nonceStr"] = nonce
		result["timestamp"] = timestamp
		result["url"] = url

		return &result, nil
	}

	return nil, nil
}

func (s *DingtalkService) GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *social.SocialUser {
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

	app := s.GetAppByAppId(appId)
	if app == nil {
		logger.Error("AppId不存在：", appId)
		return nil
	}

	if token, b := s.GetAccessToken(appId); b {
		if user, err := app.GetUserInfo(code, token); err == nil {
			if user != nil {
				var result *social.SocialUser
				if userDetail, err := app.GetUserDetail(user.Result.UserId, token); err == nil {
					mobile := userDetail.Result.Mobile
					if mobile != "" {
						// 截取mobile的最后11位数
						mobile = util.SubString(mobile, len(mobile)-11, 11)
					}

					result = &social.SocialUser{
						Name:   util.StringAddr(userDetail.Result.Name),
						Mobile: util.StringAddr(mobile),
						Email:  util.StringAddr(userDetail.Result.Email),
						Avatar: util.StringAddr(userDetail.Result.Avatar),
						OpenId: util.StringAddr(userDetail.Result.UserId),
						Title:  util.StringAddr(userDetail.Result.Title),
						Remark: util.StringAddr(userDetail.Result.Remark),
						Status: util.IntAddr(consts.STATUS_VALID),
					}
				} else {
					logger.Warn("Get user detail error: ", err.Error())
					result = &social.SocialUser{
						Name:   util.StringAddr(user.Result.Name),
						OpenId: util.StringAddr(user.Result.UserId),
						Status: util.IntAddr(consts.STATUS_VALID),
					}
				}
				// result.SetSocialId("dt", user.Result.UnionId)
				return result
			}
		} else {
			logger.Warn("User from code error: ", code, err.Error())
		}
	}

	return nil
}

type DingtalkController struct {
	DingtalkService *DingtalkService
}

func (c *DingtalkController) GetSignature(ctx *gin.Context) {
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
	if res, err := c.DingtalkService.GetSignature(ctx, appId, path, nonce, timestamp); err == nil {
		response.Success(ctx, res)
	} else {
		logger.Error("Get signature error: ", err.Error())
		response.FailMessage(ctx, 400, err.Error())
	}
}

/**
 * 处理WEB验证码的API路由
 */
func (e *DingtalkController) InitRouter(g *gin.Engine) {
	if config.Setting.Enabled {
		// 创建一个验证码路由
		wxcp := g.Group("/openapi/social/dingtalk")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			wxcp.GET("/signature", e.GetSignature) // 发送
		}
	}
}
