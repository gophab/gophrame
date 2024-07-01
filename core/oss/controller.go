package oss

import (
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/webservice/response"
)

type OssController struct {
	controller.ResourceController
	Oss OSS `inject:"oss"`
}

// 生成验证码
func (c *OssController) UploadFile(context *gin.Context) {
	_, header, err := context.Request.FormFile("file")
	if err != nil {
		response.FailMessage(context, 400, "接收文件失败")
		return
	}

	if url, _, err := c.Oss.Upload(header, "file"); err == nil {
		response.Success(context, url)
	} else {
		response.SystemErrorMessage(context, 500, err.Error())
	}
}

func (c *OssController) UploadImage(context *gin.Context) {
	_, header, err := context.Request.FormFile("file")
	if err != nil {
		response.FailMessage(context, 400, "接收文件失败")
		return
	}

	if url, _, err := c.Oss.Upload(header, "images"); err == nil {
		response.Success(context, url)
	} else {
		response.SystemErrorMessage(context, 500, err.Error())
	}
}

func (c *OssController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	// 创建一个验证码路由
	oss := g.Group("openapi/oss")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		oss.POST("/file", c.UploadFile)   //  获取验证码ID
		oss.POST("/image", c.UploadImage) //  获取验证码ID
	}
	return oss
}
