package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gophab/gophrame/core/inject"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"

	"github.com/gophab/gophrame/service"
)

type InviteCodeService struct {
	service.BaseService
	InviteCodeRepository *repository.InviteCodeRepository `inject:"inviteCodeRepository"`
	UserRepository       *repository.UserRepository       `inject:"userRepository"`
}

var inviteCodeService = &InviteCodeService{}

func init() {
	inject.InjectValue("inviteCodeService", inviteCodeService)
}

func (s *InviteCodeService) FindByInviteCode(inviteCode string) (*domain.InviteCode, error) {
	return s.InviteCodeRepository.FindByInviteCode(inviteCode)
}

func (s *InviteCodeService) GetUserInviteCode(userId string, channel string) (*domain.InviteCode, error) {
	inviteCode, err := s.InviteCodeRepository.GetUserInviteCode(userId, channel)
	if err != nil {
		return nil, err
	}

	if inviteCode == nil {
		// 实际用户
		if user, _ := s.UserRepository.GetUserById(userId); user != nil {

			if channel == "INVITE_REGISTER" {
				// 非受邀用户可以邀请注册
				if user.InviterId != nil {
					return nil, nil
				}
			}

			// create new invite code
			count := 3
			for {
				inviteCode = &domain.InviteCode{
					InviteCode:   fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)),
					UserId:       userId,
					Channel:      channel,
					InviteLimit:  0,
					InvitedLimit: 1,
				}
				res := s.InviteCodeRepository.Create(inviteCode)
				if res.Error == nil {
					break
				}

				count--
				if count < 0 {
					return nil, res.Error
				}

				time.Sleep(0)
			}
		}
	}

	return inviteCode, nil
}
