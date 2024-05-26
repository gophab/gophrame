package wxcp

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/social/wxcp/config"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type WxcpController struct {
	controller.ResourceController
	WxcpService *WxcpService `inject:"wxcpService"`
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
func (e *WxcpController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	wxcp := g.Group("/openapi/social/wxcp")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		wxcp.GET("/signature", e.GetSignature)            // 发送
		wxcp.GET("/agent-signature", e.GetAgentSignature) // 发送
	}
	return wxcp
}
