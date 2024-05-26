package wxmp

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/social/wxmp/config"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type WxmpController struct {
	controller.ResourceController
	WxmpService *WxmpService `inject:"wxmpService"`
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
func (e *WxmpController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	wxmp := g.Group("/openapi/social/wxmp")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		wxmp.GET("/signature", e.GetSignature) // 发送
	}
	return wxmp
}
