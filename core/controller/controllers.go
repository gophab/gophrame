package controller

import "github.com/gin-gonic/gin"

type Controllers struct {
	Base        string
	Handlers    []gin.HandlerFunc
	Controllers []Controller
}

// interface Controller
func (s *Controllers) InitRouter(r *gin.RouterGroup) *gin.RouterGroup {
	if s.Base != "" {
		r = r.Group(s.Base, s.Handlers...)
	}

	if len(s.Controllers) > 0 {
		// 子Controller的路由
		for _, controller := range s.Controllers {
			controller.InitRouter(r)
		}
	}

	return r
}

func (s *Controllers) AddController(c ...Controller) {
	s.Controllers = append(s.Controllers, c...)
}

func (s *Controllers) AddHandler(h ...gin.HandlerFunc) {
	s.Handlers = append(s.Handlers, h...)
}
