package service

import (
	"github.com/wjshen/gophrame/service/dto"
)

type SocialUserService interface {
	GetById(id string) (*dto.SocialUser, error)
}

func GetSocialUserService() SocialUserService {
	return _services.SocialUserService
}
