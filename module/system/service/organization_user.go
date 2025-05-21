package service

import (
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
)

type OrganizationUserService struct {
	service.BaseService
	OrganizationUserRepository *repository.OrganizationUserRepository `inject:"organizationUserRepository"`
	OrganizationRepository     *repository.OrganizationRepository     `inject:"organizationRepository"`
}

var organizationUserService *OrganizationUserService = &OrganizationUserService{}

func init() {
	inject.InjectValue("organizationUserService", organizationUserService)
}

func (u *OrganizationUserService) ListMembers(organizationId string, userName string, pageable query.Pageable) (int64, []domain.OrganizationMember) {
	return u.OrganizationUserRepository.ListMembers(organizationId, userName, pageable)
}

// 根据用户id查询所有可能的岗位节点id
func (u *OrganizationUserService) GetUserOrganizationIds(userId string) []string {
	//获取用户的所有岗位id
	organizationUsers := u.OrganizationUserRepository.GetByUserId(userId)

	organizationIds := []string{}
	for _, v := range organizationUsers {
		organizationIds = append(organizationIds, v.OrganizationId)
	}

	//根据岗位ID获取所有的岗位ID,父子级(需要去重)
	organization := u.OrganizationRepository.GetByIds(organizationIds)
	organizationIdArr := []string{}
	for _, v := range organization {
		idArr := strings.Split(v.PathInfo, ",")
		for _, vv := range idArr {
			if vv != "" {
				organizationIdArr = append(organizationIdArr, vv)
			}
		}
	}
	return organizationIdArr
}
