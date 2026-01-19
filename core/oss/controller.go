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
		response.FailMessage(context, 400, "接收图片文件失败")
		return
	}

	if url, _, err := c.Oss.Upload(header, "images"); err == nil {
		response.Success(context, url)
	} else {
		response.SystemErrorMessage(context, 500, err.Error())
	}
}

func (c *OssController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	oss := g.Group("openapi/oss")
	{
		oss.POST("/file", c.UploadFile)   //  上传文件
		oss.POST("/image", c.UploadImage) //  上传图片
	}
	return oss
}
