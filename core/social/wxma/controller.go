package wxma

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"

	"github.com/gophab/gophrame/errors"
)

type WxmaController struct {
	controller.ResourceController
	WxmaService *WxmaService `inject:"wxmaService"`
}

// 用户授权手机号绑定
func (c *WxmaController) GetPhoneNumber(ctx *gin.Context) {
	appId := request.Param(ctx, "appId").DefaultString(ctx.Request.Header.Get("X-App-Id"))
	if appId == "" {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	code, err := request.Param(ctx, "code").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.WxmaService.GetUserPhoneNumber(ctx, fmt.Sprintf("wxma:%s", appId), code)
	if err != nil {
		// 社交账号不存在
		response.SystemError(ctx, err)
		return
	}

	result = util.FullPhoneNumber(result)

	// 绑定社交账号
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	if currentUserId != "" {
		// 当前账号为社交账号，社交账号绑定手机号
		eventbus.PublishEvent("SOCIAL_BIND_MOBILE", core.M{"userId": currentUserId, "mobile": result, "social": "wxma"})
	}
	response.Success(ctx, core.M{
		"phoneNumber": result,
	})
}

type BindInfo struct {
	Code     string `json:"code"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
}

// 用户授权手机号绑定
func (c *WxmaController) BindPhoneNumber(ctx *gin.Context) {
	appId := request.Param(ctx, "appId").DefaultString(ctx.Request.Header.Get("X-App-Id"))
	if appId == "" {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	var bindInfo BindInfo
	if err := ctx.BindJSON(&bindInfo); err != nil || bindInfo.Code == "" {
		logger.Error("Bind info error: ", err.Error())
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	result, err := c.WxmaService.GetUserPhoneNumber(ctx, fmt.Sprintf("wxma:%s", appId), bindInfo.Code)
	if err != nil {
		// 社交账号不存在
		response.SystemError(ctx, err)
		return
	}

	result = util.FullPhoneNumber(result)

	// 绑定社交账号
	currentUserId := SecurityUtil.GetCurrentUserId(ctx)
	if currentUserId != "" {
		// 当前账号为社交账号，社交账号绑定手机号
		eventbus.PublishEvent("SOCIAL_BIND_MOBILE", core.M{
			"userId":   currentUserId,
			"mobile":   result,
			"nickName": bindInfo.NickName,
			"avatar":   bindInfo.Avatar,
			"social":   "wxma",
		})
	}
	response.Success(ctx, core.M{
		"phoneNumber": result,
	})
}

/**
 * 处理WEB验证码的API路由
 */
func (e *WxmaController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	wxma := g.Group("/openapi/social/wxma", security.CheckTokenVerify())
	{
		wxma.GET("/phone", e.GetPhoneNumber)   // 发送
		wxma.POST("/phone", e.BindPhoneNumber) // 发送
	}
	return wxma
}
