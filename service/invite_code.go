package service

import "github.com/gophab/gophrame/service/dto"

type InviteCodeService interface {
	FindByInviteCode(code string) (*dto.InviteCode, error)
}

func GetInviteCodeService() InviteCodeService {
	return _services.InviteCodeService
}
