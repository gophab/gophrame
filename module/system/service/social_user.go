package service

import (
	"strings"
	"time"

	"github.com/gophab/gophrame/core/consts"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
)

type SocialUserService struct {
	service.BaseService
	SocialUserRepository *repository.SocialUserRepository `inject:"socialUserRepository"`
}

var socialUserService *SocialUserService = &SocialUserService{}

func init() {
	eventbus.RegisterEventListener("USER_LOGIN", socialUserService.onUserLogin)
	inject.InjectValue("socialUserService", socialUserService)
}

func GetSocialUserService() *SocialUserService {
	return socialUserService
}

func (s *SocialUserService) GetById(id string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetById(id)
}

func (s *SocialUserService) GetBySocialId(socialType string, socialId string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetBySocialId(socialType, socialId)
}

func (s *SocialUserService) GetByUserId(socialType string, userId string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetByUserId(socialType, userId)
}

func (s *SocialUserService) CreateSocialUser(socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	socialUser.Status = util.IntAddr(consts.STATUS_VALID)
	if res := s.SocialUserRepository.Create(socialUser); res.Error == nil && res.RowsAffected > 0 {
		return socialUser, nil
	} else {
		return nil, res.Error
	}
}

func (s *SocialUserService) UpdateSocialUser(socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	exists, err := s.GetById(socialUser.Id)
	if err != nil || exists == nil {
		return nil, err
	}

	updated := false

	if socialUser.Name != nil {
		exists.Name = socialUser.Name
		updated = true
	}
	if socialUser.NickName != nil {
		exists.NickName = socialUser.NickName
		updated = true
	}
	if socialUser.Avatar != nil {
		exists.Avatar = socialUser.Avatar
		updated = true
	}
	if socialUser.Mobile != nil {
		exists.Mobile = socialUser.Mobile
		updated = true
	}
	if socialUser.Email != nil {
		exists.Email = socialUser.Email
		updated = true
	}
	if socialUser.Title != nil {
		exists.Title = socialUser.Title
		updated = true
	}
	if socialUser.Remark != nil {
		exists.Remark = socialUser.Remark
		updated = true
	}
	if socialUser.OpenId != nil {
		exists.OpenId = socialUser.OpenId
		updated = true
	}
	if socialUser.SocialId != nil {
		exists.SocialId = socialUser.SocialId
		updated = true
	}
	if socialUser.UserId != nil {
		exists.UserId = socialUser.UserId
		updated = true
	}
	if socialUser.Status != nil {
		exists.Status = socialUser.Status
		updated = true
	}

	if updated {
		if res := s.SocialUserRepository.Omit("type", "login_times", "last_login_time", "last_login_ip", "created_time", "last_modified_time").Save(exists); res.Error == nil {
			return exists, nil
		} else {
			return nil, res.Error
		}
	}
	return exists, nil
}

func (s *SocialUserService) BoundSocialUser(socialUserId string, userId string, socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	exists, err := s.GetById(socialUserId)
	if err != nil || exists == nil {
		return nil, err
	}

	exists.UserId = &userId

	if socialUser != nil {
		if socialUser.Name != nil {
			exists.Name = socialUser.Name
		}
		if socialUser.NickName != nil {
			exists.NickName = socialUser.NickName
		}
		if socialUser.Avatar != nil {
			exists.Avatar = socialUser.Avatar
		}
		if socialUser.Mobile != nil {
			exists.Mobile = socialUser.Mobile
		}
		if socialUser.Email != nil {
			exists.Email = socialUser.Email
		}
		if socialUser.Title != nil {
			exists.Title = socialUser.Title
		}
		if socialUser.Remark != nil {
			exists.Remark = socialUser.Remark
		}
	}

	if res := s.SocialUserRepository.Omit("type", "login_times", "last_login_time", "last_login_ip", "created_time", "last_modified_time").Save(exists); res.Error == nil {
		if exists.OpenId != nil && exists.Id != socialUser.Type+"_"+*socialUser.OpenId {
			s.BoundSocialUser(socialUser.Type+"_"+*socialUser.OpenId, userId, socialUser)
		}
		return exists, nil
	} else {
		return nil, res.Error
	}
}

func (s *SocialUserService) onUserLogin(event string, args ...interface{}) {
	userId := ""
	data := map[string]string{}

	if len(args) > 0 {
		for i, v := range args {
			if i == 0 {
				userId = v.(string)
			}
			if i == 1 {
				data = v.(map[string]string)
			}
		}
	}

	if strings.HasPrefix(userId, "sns:") {
		userId, _ := strings.CutPrefix(userId, "sns:")
		if socialUser, err := s.GetById(userId); err != nil || socialUser == nil {
			return
		} else {
			if socialUser.UserId != nil {
				eventbus.PublishEvent("USER_LOGIN", socialUser.UserId, data)
			}

			socialUser.LastLoginTime = util.TimeAddr(time.Now())
			socialUser.LastLoginIp = util.StringAddr(data["IP"])
			socialUser.LoginTimes = socialUser.LoginTimes + 1
			if res := s.SocialUserRepository.Save(socialUser); res.Error != nil {
				logger.Error("Save SocialUser Error: ", res.Error.Error())
				return
			}
		}
	}
}
