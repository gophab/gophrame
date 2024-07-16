package controller

import "github.com/gin-gonic/gin"

type ResourceHandler struct {
	HttpMethod   string
	ResourcePath string
	Handler      gin.HandlerFunc
}

type Controller interface {
	InitRouter(*gin.RouterGroup) *gin.RouterGroup
}

type ResourceController struct {
	Controllers
	ResourceHandlers []ResourceHandler
}

func (c *ResourceController) RegisterRoute(httpMethod string, path string, handler gin.HandlerFunc) {
	c.ResourceHandlers = append(c.ResourceHandlers, ResourceHandler{HttpMethod: httpMethod, ResourcePath: path, Handler: handler})
}

func (c *ResourceController) RegisterResourceHandler(h ...ResourceHandler) {
	c.ResourceHandlers = append(c.ResourceHandlers, h...)
}

func (c *ResourceController) SetResourceHandlers(handlers []ResourceHandler) {
	c.ResourceHandlers = handlers
}

func (c *ResourceController) InitRouter(r *gin.RouterGroup) *gin.RouterGroup {
	// 如果有子路由的话
	r = c.Controllers.InitRouter(r)
	if len(c.ResourceHandlers) > 0 {
		for _, handler := range c.ResourceHandlers {
			r.Handle(handler.HttpMethod, handler.ResourcePath, SetGlobalContext(), handler.Handler)
		}
	}
	return r
}
