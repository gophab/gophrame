package dingtalk

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/social/dingtalk/config"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type DingtalkController struct {
	controller.ResourceController
	DingtalkService *DingtalkService `inject:"dingtalkService"`
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
func (e *DingtalkController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	dingtalk := g.Group("/openapi/social/dingtalk")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		dingtalk.GET("/signature", e.GetSignature) // 发送
	}
	return dingtalk
}
