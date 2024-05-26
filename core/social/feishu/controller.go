package feishu

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/security"
	"github.com/gophab/gophrame/core/social/feishu/config"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type FeishuController struct {
	controller.ResourceController
	FeishuService *FeishuService `inject:"feishuService"`
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
func (e *FeishuController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	feishu := g.Group("/openapi/social/feishu")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		feishu.GET("/signature", e.GetSignature)                                  // 发送
		feishu.POST("/speech-text", security.HandleTokenVerify(), e.SpeechToText) // 发送
	}
	return feishu
}
