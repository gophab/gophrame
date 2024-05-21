package mapi

import (
	"github.com/wjshen/gophrame/core/controller"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/webservice/request"
	"github.com/wjshen/gophrame/core/webservice/response"
	"github.com/wjshen/gophrame/service"

	"github.com/gin-gonic/gin"
)

type OrganizationUserMController struct {
	controller.ResourceController
	OrganizationUserService *service.OrganizationUserService `inject:"organizationUserService"`
}

var organizationUserMController *OrganizationUserMController = &OrganizationUserMController{}

func init() {
	inject.InjectValue("organizationUserMController", organizationUserMController)
}

// 成员
func (m *OrganizationUserMController) AfterInitialize() {
	m.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/organization/:id/members", Handler: m.GetOrganizationMemebers},
	})
}

func (o *OrganizationUserMController) GetOrganizationMemebers(c *gin.Context) {
	organizationId := request.Param(c, "organizationId").Int64()
	name := request.Param(c, "name").DefaultString("")
	pageable := query.GetPageable(c)

	count, list := o.OrganizationUserService.ListMembers(organizationId, name, pageable)
	response.Page(c, count, list)
}
